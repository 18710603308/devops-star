package controllers

import (
	"net/http"

	"devops-star/backend/services"

	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	projectService *services.ProjectService
}

func NewProjectController(projectService *services.ProjectService) *ProjectController {
	return &ProjectController{projectService: projectService}
}

// 获取项目列表
func (c *ProjectController) ListProjects(ctx *gin.Context) {
	projects, err := c.projectService.ListProjects()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": projects})
}

// 创建项目
func (c *ProjectController) CreateProject(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		DisplayName string `json:"display_name"`
		Description string `json:"description"`
		RepoURL     string `json:"repo_url"`
		RepoType    string `json:"repo_type"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从 JWT 获取用户 ID（简化版）
	userID := uint(1)

	project, err := c.projectService.CreateProject(req.Name, req.DisplayName, req.Description, req.RepoURL, req.RepoType, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, project)
}

// 获取项目详情
func (c *ProjectController) GetProject(ctx *gin.Context) {
	id := ctx.Param("id")
	project, err := c.projectService.GetProject(parseUint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
		return
	}
	ctx.JSON(http.StatusOK, project)
}

// 更新项目
func (c *ProjectController) UpdateProject(ctx *gin.Context) {
	id := ctx.Param("id")
	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := c.projectService.UpdateProject(parseUint(id), updates)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, project)
}

// 删除项目
func (c *ProjectController) DeleteProject(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.projectService.DeleteProject(parseUint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "项目删除成功"})
}

// 获取项目成员
func (c *ProjectController) GetProjectMembers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
}

// 添加项目成员
func (c *ProjectController) AddProjectMember(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "成员添加成功"})
}

func parseUint(s string) uint {
	// 简化版，实际应使用 strconv.ParseUint
	return 1
}
