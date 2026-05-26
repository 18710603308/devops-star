package services

import (
	"devops-star/backend/config"
	"devops-star/backend/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 创建测试数据库连接（使用 SQLite 内存数据库）
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("迁移数据库失败: %v", err)
	}

	return db
}

// TestAuthService_Register 测试用户注册
func TestAuthService_Register(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	service := NewAuthService(db, cfg)

	tests := []struct {
		name     string
		username string
		password string
		email    string
		wantErr  bool
	}{
		{
			name:     "有效注册",
			username: "testuser",
			password: "password123",
			email:    "test@example.com",
			wantErr:  false,
		},
		{
			name:     "用户名已存在",
			username: "testuser", // 重复
			password: "password456",
			email:    "test2@example.com",
			wantErr:  true,
		},
		{
			name:     "邮箱已存在",
			username: "testuser2",
			password: "password789",
			email:    "test@example.com", // 重复
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.Register(tt.username, tt.password, tt.email)

			if tt.wantErr {
				assert.Error(t, err, "期望返回错误")
			} else {
				assert.NoError(t, err, "不期望返回错误")
				assert.NotNil(t, user, "用户不应为 nil")
				assert.Equal(t, tt.username, user.Username, "用户名不匹配")
				assert.Equal(t, tt.email, user.Email, "邮箱不匹配")
				assert.NotEmpty(t, user.Password, "密码不应为空")
				// 密码应该被加密（不是原始密码）
				assert.NotEqual(t, tt.password, user.Password, "密码应该被加密")
			}
		})
	}
}

// TestAuthService_Login 测试用户登录
func TestAuthService_Login(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	service := NewAuthService(db, cfg)

	// 先注册一个用户
	_, err := service.Register("logintest", "password123", "login@test.com")
	assert.NoError(t, err, "注册用户应该成功")

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "有效登录（使用用户名）",
			username: "logintest",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "有效登录（使用邮箱）",
			username: "login@test.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "密码错误",
			username: "logintest",
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "用户不存在",
			username: "nonexistent",
			password: "password123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, token, err := service.Login(tt.username, tt.password)

			if tt.wantErr {
				assert.Error(t, err, "期望返回错误")
				assert.Empty(t, token, "token 应该为空")
			} else {
				assert.NoError(t, err, "不期望返回错误")
				assert.NotNil(t, user, "用户不应为 nil")
				assert.NotEmpty(t, token, "token 不应为空")
				// 验证 token 格式（应该是 3 段，用 . 分隔）
				assert.Contains(t, token, ".", "token 应该包含 . 分隔符")
			}
		})
	}
}

// TestGenerateToken 测试 JWT Token 生成
func TestGenerateToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}
	service := NewAuthService(nil, cfg)

	userID := uint(1)
	username := "testuser"

	token := generateToken(userID, username, cfg.JWTSecret)

	assert.NotEmpty(t, token, "token 不应为空")
	assert.Contains(t, token, ".", "token 应该包含 . 分隔符")

	// 验证 token 可以被解析
	parsedUserID, parsedUsername, err := ValidateToken(token, cfg.JWTSecret)
	assert.NoError(t, err, "解析 token 不应该出错")
	assert.Equal(t, userID, parsedUserID, "用户 ID 不匹配")
	assert.Equal(t, username, parsedUsername, "用户名不匹配")
}

// TestValidateToken 测试 JWT Token 验证
func TestValidateToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	tests := []struct {
		name        string
		tokenString string
		secret      string
		wantErr     bool
	}{
		{
			name:        "有效 token",
			tokenString: generateToken(1, "testuser", cfg.JWTSecret),
			secret:      cfg.JWTSecret,
			wantErr:     false,
		},
		{
			name:        "错误的安全密钥",
			tokenString: generateToken(1, "testuser", cfg.JWTSecret),
			secret:      "wrong-secret",
			wantErr:     true,
		},
		{
			name:        "无效的 token 格式",
			tokenString: "invalid.token.format",
			secret:      cfg.JWTSecret,
			wantErr:     true,
		},
		{
			name:        "空的 token",
			tokenString: "",
			secret:      cfg.JWTSecret,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, username, err := ValidateToken(tt.tokenString, tt.secret)

			if tt.wantErr {
				assert.Error(t, err, "期望返回错误")
			} else {
				assert.NoError(t, err, "不期望返回错误")
				assert.NotZero(t, userID, "用户 ID 不应为 0")
				assert.NotEmpty(t, username, "用户名不应为空")
			}
		})
	}
}

// TestPasswordHashing 测试密码加密和验证
func TestPasswordHashing(t *testing.T) {
	password := "mySecurePassword123!"

	// 加密密码
	hashedPassword, err := hashPassword(password)
	assert.NoError(t, err, "加密密码不应该出错")
	assert.NotEmpty(t, hashedPassword, "加密后的密码不应为空")
	assert.NotEqual(t, password, hashedPassword, "原始密码和加密后的密码应该不同")

	// 验证正确密码
	err = checkPassword(hashedPassword, password)
	assert.NoError(t, err, "验证正确密码不应该出错")

	// 验证错误密码
	err = checkPassword(hashedPassword, "wrongPassword")
	assert.Error(t, err, "验证错误密码应该返回错误")
}
