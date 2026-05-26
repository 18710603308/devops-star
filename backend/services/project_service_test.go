package services

import (
	"devops-star/backend/config"
	"devops-star/backend/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupProjectTestDB 创建项目测试数据库
func setupProjectTestDB(t *testing.T) (*gorm.DB, *config.Config) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&models.Project{})
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

// TestProjectService_CreateProject 测试创建项目
func TestProjectService_CreateProject(t *testing.T) {
	db, cfg := setupProjectTestDB(t)
	service := NewProjectService(db, cfg, nil)

	tests := []struct {
		name        string
		projectName string
		displayName string
		description string
		repoURL     string
		repoType    string
		createdBy   uint
		wantErr     bool
	}{
		{
			name:        "有效项目创建",
			projectName: "test-project",
			displayName: "测试项目",
			description: "这是一个测试项目",
			repoURL:     "",
			repoType:    "git",
			createdBy:   1,
			wantErr:     false,
		},
		{
			name:        "项目名称包含大写（应该被转为小写）",
			projectName: "MyProject",
			displayName: "我的项目",
			description: "",
			repoURL:     "",
			repoType:    "git",
			createdBy:   1,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project, err := service.CreateProject(tt.projectName, tt.displayName, tt.description, tt.repoURL, tt.repoType, tt.createdBy)

			if tt.wantErr {
				assert.Error(t, err, "期望返回错误")
			} else {
				assert.NoError(t, err, "不期望返回错误")
				assert.NotNil(t, project, "项目不应为 nil")
				assert.Equal(t, tt.projectName, project.Name, "项目名称不匹配")
				assert.Equal(t, tt.displayName, project.DisplayName, "项目显示名称不匹配")
				// 验证仓库名称被格式化为小写
				assert.Equal(t, "test-project", project.Name, "仓库名称应该被转为小写")
			}
		})
	}
}

// TestProjectService_GetProject 测试获取项目
func TestProjectService_GetProject(t *testing.T) {
	db, cfg := setupProjectTestDB(t)
	service := NewProjectService(db, cfg, nil)

	// 先创建一个项目
	project, err := service.CreateProject("get-test", "获取测试", "测试描述", "", "git", 1)
	assert.NoError(t, err, "创建项目应该成功")
	assert.NotNil(t, project, "创建的项目不应为 nil")

	// 测试获取存在的项目
	gotProject, err := service.GetProject(project.ID)
	assert.NoError(t, err, "获取项目不应该出错")
	assert.NotNil(t, gotProject, "获取的项目不应为 nil")
	assert.Equal(t, project.ID, gotProject.ID, "项目 ID 不匹配")
	assert.Equal(t, project.Name, gotProject.Name, "项目名称不匹配")

	// 测试获取不存在的项目
	_, err = service.GetProject(999)
	assert.Error(t, err, "获取不存在的项目应该返回错误")
}

// TestProjectService_UpdateProject 测试更新项目
func TestProjectService_UpdateProject(t *testing.T) {
	db, cfg := setupProjectTestDB(t)
	service := NewProjectService(db, cfg, nil)

	// 先创建一个项目
	project, err := service.CreateProject("update-test", "更新测试", "原始描述", "", "git", 1)
	assert.NoError(t, err, "创建项目应该成功")

	// 更新项目
	newDisplayName := "更新后的名称"
	newDescription := "更新后的描述"
	updatedProject, err := service.UpdateProject(project.ID, newDisplayName, newDescription)
	assert.NoError(t, err, "更新项目不应该出错")
	assert.NotNil(t, updatedProject, "更新的项目不应为 nil")
	assert.Equal(t, newDisplayName, updatedProject.DisplayName, "更新后的显示名称不匹配")
	assert.Equal(t, newDescription, updatedProject.Description, "更新后的描述不匹配")
}

// TestProjectService_DeleteProject 测试删除项目
func TestProjectService_DeleteProject(t *testing.T) {
	db, cfg := setupProjectTestDB(t)
	service := NewProjectService(db, cfg, nil)

	// 先创建一个项目
	project, err := service.CreateProject("delete-test", "删除测试", "测试描述", "", "git", 1)
	assert.NoError(t, err, "创建项目应该成功")

	// 删除项目
	err = service.DeleteProject(project.ID)
	assert.NoError(t, err, "删除项目不应该出错")

	// 验证项目已被删除（软删除）
	_, err = service.GetProject(project.ID)
	assert.Error(t, err, "删除后获取项目应该返回错误")
}

// TestProjectService_ListProjects 测试列出项目
func TestProjectService_ListProjects(t *testing.T) {
	db, cfg := setupProjectTestDB(t)
	service := NewProjectService(db, cfg, nil)

	// 创建多个项目
	service.CreateProject("list-test-1", "列表测试1", "", "", "git", 1)
	service.CreateProject("list-test-2", "列表测试2", "", "", "git", 1)
	service.CreateProject("list-test-3", "列表测试3", "", "", "git", 1)

	// 列出所有项目
	projects, err := service.ListProjects(0)
	assert.NoError(t, err, "列出项目不应该出错")
	assert.GreaterOrEqual(t, len(projects), 3, "应该至少有 3 个项目")
}
