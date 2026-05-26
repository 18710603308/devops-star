package services

import (
	"devops-star/backend/config"
	"devops-star/backend/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{DB: db, Cfg: cfg}
}

// 用户注册
func (s *AuthService) Register(username, password, email string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Role:     "user",
		Active:   true,
	}

	if err := s.DB.Create(user).Error; err != nil {
		return nil, err
	}

	user.Password = "" // 清除密码
	return user, nil
}

// 用户登录
func (s *AuthService) Login(username, password string) (*models.User, string, error) {
	var user models.User
	if err := s.DB.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", err
	}

	// 生成 JWT Token
	token := generateToken(user.ID, user.Username, s.Cfg.JWTSecret)

	user.Password = ""
	return &user, token, nil
}

// 获取用户信息
func (s *AuthService) GetUserInfo(userID uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	user.Password = ""
	return &user, nil
}

// 生成 JWT Token（使用 golang-jwt/jwt/v5）
func generateToken(userID uint, username string, secret string) string {
	// 创建 JWT Claims
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整 token 字符串
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		// 降级处理：返回简化 token
		return "error-generating-token"
	}

	return tokenString
}

// 验证 JWT Token（导出函数）
func ValidateToken(tokenString string, secret string) (uint, string, error) {
	// 解析并验证 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, "", err
	}

	// 提取 claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		username := claims["username"].(string)
		return userID, username, nil
	}

	return 0, "", fmt.Errorf("invalid token")
}
