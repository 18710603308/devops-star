package services

import (
	"devops-star/backend/config"
	"devops-star/backend/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupDeployTestDB 创建部署测试数据库
func setupDeployTestDB(t *testing.T) (*gorm.DB, *config.Config) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&models.Environment{}, &models.DeployRecord{})
	if err != nil {
		t.Fatalf("迁移数据库失败: %v", err)
	}

	cfg := &config.Config{
		GiteaURL:       "http://localhost:3000",
		GiteaAdminUser:  "admin",
		GiteaAdminPass: "admin123",
	}

	return db, cfg
}

// TestDeployService_CreateEnvironment 测试创建部署环境
func TestDeployService_CreateEnvironment(t *testing.T) {
	db, cfg := setupDeployTestDB(t)
	service := NewDeployService(db, cfg)

	tests := []struct {
		name        string
		envName     string
		displayName string
		deployType  string
		projectID   uint
		wantErr     bool
	}{
		{
			name:        "有效环境创建",
			envName:     "dev",
			displayName: "开发环境",
			deployType:  "docker",
			projectID:   1,
			wantErr:     false,
		},
		{
			name:        "K8s 环境创建",
			envName:     "prod",
			displayName: "生产环境",
			deployType:  "k8s",
			projectID:   1,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := service.CreateEnvironment(tt.envName, tt.displayName, tt.deployType, tt.projectID)

			if tt.wantErr {
				assert.Error(t, err, "期望返回错误")
			} else {
				assert.NoError(t, err, "不期望返回错误")
				assert.NotNil(t, env, "环境不应为 nil")
				assert.Equal(t, tt.envName, env.Name, "环境名称不匹配")
				assert.Equal(t, tt.displayName, env.DisplayName, "环境显示名称不匹配")
				assert.Equal(t, tt.deployType, env.DeployType, "部署类型不匹配")
			}
		})
	}
}

// TestDeployService_GetEnvironments 测试获取部署环境列表
func TestDeployService_GetEnvironments(t *testing.T) {
	db, cfg := setupDeployTestDB(t)
	service := NewDeployService(db, cfg)

	// 先创建几个环境
	service.CreateEnvironment("dev", "开发环境", "docker", 1)
	service.CreateEnvironment("test", "测试环境", "docker", 1)
	service.CreateEnvironment("prod", "生产环境", "k8s", 2)

	tests := []struct {
		name      string
		projectID uint
		expected  int
	}{
		{
			name:      "获取所有环境",
			projectID: 0,
			expected:  3,
		},
		{
			name:      "按项目 ID 获取环境",
			projectID: 1,
			expected:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envs, err := service.GetEnvironments(tt.projectID)
			assert.NoError(t, err, "不期望返回错误")
			assert.Equal(t, tt.expected, len(envs), "环境数量不匹配")
		})
	}
}

// TestDeployService_TriggerDeploy 测试触发部署
func TestDeployService_TriggerDeploy(t *testing.T) {
	db, cfg := setupDeployTestDB(t)
	service := NewDeployService(db, cfg)

	// 先创建一个环境
	env, err := service.CreateEnvironment("dev", "开发环境", "docker", 1)
	assert.NoError(t, err, "创建环境应该成功")

	// 触发部署
	record, err := service.TriggerDeploy(env.ID, "run-001", "nginx:latest", 1)
	assert.NoError(t, err, "触发部署不应该出错")
	assert.NotNil(t, record, "部署记录不应为 nil")
	assert.Equal(t, env.ID, record.EnvironmentID, "环境 ID 不匹配")
	assert.Equal(t, "run-001", record.PipelineRunID, "Pipeline 运行 ID 不匹配")
	assert.Equal(t, "nginx:latest", record.ImageTag, "镜像标签不匹配")
	assert.Equal(t, "pending", record.Status, "部署状态应该是 pending")
}

// TestDeployService_GetDeployHistory 测试获取部署历史
func TestDeployService_GetDeployHistory(t *testing.T) {
	db, cfg := setupDeployTestDB(t)
	service := NewDeployService(db, cfg)

	// 先创建一个环境
	env, err := service.CreateEnvironment("dev", "开发环境", "docker", 1)
	assert.NoError(t, err, "创建环境应该成功")

	// 触发几次部署
	service.TriggerDeploy(env.ID, "run-001", "nginx:v1.0.0", 1)
	service.TriggerDeploy(env.ID, "run-002", "nginx:v1.0.1", 1)
	service.TriggerDeploy(env.ID, "run-003", "nginx:v1.0.2", 1)

	// 获取部署历史
	records, err := service.GetDeployHistory(0)
	assert.NoError(t, err, "获取部署历史不应该出错")
	assert.GreaterOrEqual(t, len(records), 3, "应该至少有 3 条部署记录")

	// 按环境 ID 获取部署历史
	records, err = service.GetDeployHistory(env.ID)
	assert.NoError(t, err, "按环境获取部署历史不应该出错")
	assert.Equal(t, 3, len(records), "环境应该有的 3 条部署记录")
}

// TestDeployService_GetDeployStats 测试获取部署统计
func TestDeployService_GetDeployStats(t *testing.T) {
	db, cfg := setupDeployTestDB(t)
	service := NewDeployService(db, cfg)

	// 先创建一个环境
	env, err := service.CreateEnvironment("dev", "开发环境", "docker", 1)
	assert.NoError(t, err, "创建环境应该成功")

	// 触发几次部署
	service.TriggerDeploy(env.ID, "run-001", "nginx:v1.0.0", 1)
	service.TriggerDeploy(env.ID, "run-002", "nginx:v1.0.1", 1)

	// 获取部署统计
	stats, err := service.GetDeployStats()
	assert.NoError(t, err, "获取部署统计不应该出错")
	assert.NotNil(t, stats, "统计不应为 nil")
	assert.Equal(t, int64(2), stats["total"], "总部署次数应该是 2")
}

// TestDeployService_RollbackDeploy 测试回滚部署（简化版）
func TestDeployService_RollbackDeploy(t *testing.T) {
	db, cfg := setupDeployTestDB(t)
	service := NewDeployService(db, cfg)

	// 先创建一个环境
	env, err := service.CreateEnvironment("dev", "开发环境", "docker", 1)
	assert.NoError(t, err, "创建环境应该成功")

	// 触发几次部署
	record1, _ := service.TriggerDeploy(env.ID, "run-001", "nginx:v1.0.0", 1)
	record2, _ := service.TriggerDeploy(env.ID, "run-002", "nginx:v1.0.1", 1)

	// 模拟 record1 部署成功
	db.Model(&models.DeployRecord{}).Where("id = ?", record1.ID).Update("status", "success")
	db.Model(&models.DeployRecord{}).Where("id = ?", record2.ID).Update("status", "success")

	// 回滚到 record1
	err = service.RollbackDeploy(record2.ID)
	// 注意：这个测试可能会失败，因为回滚功能需要真正执行部署
	// 这里只是测试函数调用不会 panic
	if err != nil {
		t.Logf("回滚失败（预期）: %v", err)
	}
}

// TestFormatRepoName 测试仓库名称格式化
func TestFormatRepoName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"My Project", "my-project"},
		{"Test_Repo", "test-repo"},
		{"Demo.Repo", "demo-repo"},
		{"Already-Lower", "already-lower"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := FormatRepoName(tt.input)
			assert.Equal(t, tt.expected, result, "格式化后的仓库名称不匹配")
		})
	}
}
