package services

import (
	"devops-star/backend/config"
	"devops-star/backend/models"
	"fmt"

	"gorm.io/gorm"
)

type ProjectService struct {
	DB           *gorm.DB
	Cfg          *config.Config
	GiteaService *GiteaService
}

func NewProjectService(db *gorm.DB, cfg *config.Config, giteaService *GiteaService) *ProjectService {
	return &ProjectService{DB: db, Cfg: cfg, GiteaService: giteaService}
}

// 创建项目
func (s *ProjectService) CreateProject(name, displayName, description, repoURL, repoType string, createdBy uint) (*models.Project, error) {
	// 格式化仓库名称（Gitea 要求小写）
	repoName := FormatRepoName(name)

	// 在 Gitea 中创建仓库
	if s.GiteaService != nil {
		fmt.Printf("🔧 在 Gitea 中创建仓库: %s\n", repoName)
		repo, err := s.GiteaService.CreateRepo(repoName, description, false)
		if err != nil {
			// 仓库创建失败，但不阻塞项目创建（可能是已存在）
			fmt.Printf("⚠️  Gitea 仓库创建失败（可能已存在）: %v\n", err)
			// 继续执行，使用用户提供的 repoURL 或生成默认 URL
		} else {
			// 仓库创建成功，使用 Gitea 提供的克隆地址
			repoURL = repo.CloneURL
			fmt.Printf("✅ Gitea 仓库创建成功: %s\n", repo.HTMLURL)

			// 自动提交示例 workflow 文件到仓库
			fmt.Printf("📦 提交示例 workflow 文件到仓库...\n")
			workflowContent := `name: DevOpsStar CI/CD 示例

on:
  push:
    branches: [main, master, develop]
  pull_request:
    branches: [main, master]
  workflow_dispatch:  # 允许手动触发

jobs:
  build:
    name: 构建项目
    runs-on: docker
    container:
      image: node:20-alpine
    steps:
      - name: 检出代码
        uses: actions/checkout@v4

      - name: 设置 Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: 安装依赖
        run: |
          echo "📦 安装依赖..."
          npm ci --registry=https://registry.npmmirror.com
          echo "✅ 依赖安装完成"

      - name: 执行构建
        run: |
          echo "🔧 开始构建..."
          npm run build
          echo "✅ 构建完成"

      - name: 上传构建产物
        uses: actions/upload-artifact@v4
        with:
          name: build-artifacts
          path: dist/
          retention-days: 7
`

			commitErr := s.GiteaService.CommitFile(repoName, ".gitea/workflows/ci.yml", workflowContent, "Initialize CI/CD pipeline")
			if commitErr != nil {
				fmt.Printf("⚠️  提交 workflow 文件失败: %v\n", commitErr)
			} else {
				fmt.Printf("✅ 示例 workflow 文件已提交到仓库\n")
			}
		}
	}

	// 如果 repoURL 仍为空，生成默认 URL
	if repoURL == "" && s.Cfg.GiteaURL != "" {
		repoURL = fmt.Sprintf("%s/%s/%s.git", s.Cfg.GiteaURL, s.Cfg.GiteaAdminUser, repoName)
	}

	project := &models.Project{
		Name:        repoName,
		DisplayName: displayName,
		Description: description,
		RepoURL:     repoURL,
		RepoType:    repoType,
		CreatedBy:   createdBy,
	}

	if err := s.DB.Create(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

// 获取项目列表
func (s *ProjectService) ListProjects() ([]models.Project, error) {
	var projects []models.Project
	if err := s.DB.Where("deleted_at IS NULL").Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// 获取单个项目
func (s *ProjectService) GetProject(id uint) (*models.Project, error) {
	var project models.Project
	if err := s.DB.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// 更新项目
func (s *ProjectService) UpdateProject(id uint, updates map[string]interface{}) (*models.Project, error) {
	var project models.Project
	if err := s.DB.First(&project, id).Error; err != nil {
		return nil, err
	}
	if err := s.DB.Model(&project).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.DB.First(&project, id)
	return &project, nil
}

// 删除项目
func (s *ProjectService) DeleteProject(id uint) error {
	return s.DB.Delete(&models.Project{}, id).Error
}
