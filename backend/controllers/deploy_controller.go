package controllers

import (
	"fmt"
	"net/http"

	"devops-star/backend/services"

	"github.com/gin-gonic/gin"
)

type DeployController struct {
	deployService  *services.DeployService
	notifyService *services.NotificationService
}

func NewDeployController(deployService *services.DeployService, notifyService *services.NotificationService) *DeployController {
	return &DeployController{
		deployService:  deployService,
		notifyService: notifyService,
	}
}

// 获取环境列表
func (c *DeployController) GetEnvironments(ctx *gin.Context) {
	projectID := uint(0)
	if ctx.Query("project_id") != "" {
		// 简化：实际应从 JWT 或参数获取
		projectID = 0 // 0 表示获取所有
	}

	envs, err := c.deployService.GetEnvironments(projectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": envs})
}

// 创建环境
func (c *DeployController) CreateEnvironment(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		DisplayName string `json:"display_name"`
		DeployType  string `json:"deploy_type" binding:"required"`
		ProjectID   uint   `json:"project_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	env, err := c.deployService.CreateEnvironment(req.Name, req.DisplayName, req.DeployType, req.ProjectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, env)
}

// 触发部署
func (c *DeployController) TriggerDeploy(ctx *gin.Context) {
	var req struct {
		EnvironmentID uint   `json:"environment_id" binding:"required"`
		PipelineRunID string `json:"pipeline_run_id"`
		ImageTag      string `json:"image_tag"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从 JWT 获取用户信息（简化版）
	deployedBy := uint(1)

	record, err := c.deployService.TriggerDeploy(req.EnvironmentID, req.PipelineRunID, req.ImageTag, deployedBy)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 发送通知
	go c.notifyService.SendNotificationAsync("部署任务已创建: "+record.ImageTag, "always")

	ctx.JSON(http.StatusOK, gin.H{"message": "部署任务已创建", "deploy_id": record.ID})
}

// 获取部署历史
func (c *DeployController) GetDeployHistory(ctx *gin.Context) {
	envID := uint(0)
	if ctx.Query("environment_id") != "" {
		// 简化：解析 environment_id
	}

	records, err := c.deployService.GetDeployHistory(envID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": records})
}

// 回滚部署
func (c *DeployController) RollbackDeploy(ctx *gin.Context) {
	recordID := ctx.Param("id")

	var req struct {
		TargetRecordID uint `json:"target_record_id"` // 可选：回滚到指定记录
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 简化：使用 URL 中的 recordID
	rid := uint(0)
	fmt.Sscanf(recordID, "%d", &rid)

	if rid == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的部署记录 ID"})
		return
	}

	if err := c.deployService.RollbackDeploy(rid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "回滚任务已创建"})
}

// 获取部署日志
func (c *DeployController) GetDeployLogs(ctx *gin.Context) {
	recordID := ctx.Param("id")

	rid := uint(0)
	fmt.Sscanf(recordID, "%d", &rid)

	if rid == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的部署记录 ID"})
		return
	}

	logs, err := c.deployService.GetDeployLogs(rid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"logs": logs})
}

// 检查部署健康状态
func (c *DeployController) CheckHealth(ctx *gin.Context) {
	envID := ctx.Param("id")

	eid := uint(0)
	fmt.Sscanf(envID, "%d", &eid)

	if eid == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的环境 ID"})
		return
	}

	status, err := c.deployService.CheckDeploymentHealth(eid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": status})
}
