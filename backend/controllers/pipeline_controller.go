package controllers

import (
	"net/http"

	"devops-star/backend/services"

	"github.com/gin-gonic/gin"
)

type PipelineController struct {
	pipelineService  *services.PipelineService
	notifyService    *services.NotificationService
}

func NewPipelineController(pipelineService *services.PipelineService, notifyService *services.NotificationService) *PipelineController {
	return &PipelineController{pipelineService: pipelineService, notifyService: notifyService}
}

// 获取流水线列表
func (c *PipelineController) ListPipelines(ctx *gin.Context) {
	projectID := parseUint(ctx.Query("project_id"))
	pipelines, err := c.pipelineService.ListPipelines(projectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": pipelines})
}

// 创建流水线
func (c *PipelineController) CreatePipeline(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		ConfigYAML  string `json:"config_yaml"`
		ProjectID   uint   `json:"project_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从 JWT 获取用户 ID（简化版）
	userID := uint(1)

	pipeline, err := c.pipelineService.CreatePipeline(req.Name, req.Description, req.ConfigYAML, req.ProjectID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, pipeline)
}

// 获取流水线详情
func (c *PipelineController) GetPipeline(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	pipeline, err := c.pipelineService.GetPipeline(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "流水线不存在"})
		return
	}
	ctx.JSON(http.StatusOK, pipeline)
}

// 触发流水线
func (c *PipelineController) TriggerPipeline(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	// 从 JWT 获取用户名（简化版）
	trigger := "admin"
	branch := ctx.Query("branch")
	if branch == "" {
		branch = "main"
	}

	run, err := c.pipelineService.TriggerPipeline(id, trigger, branch)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取流水线信息用于通知
	pipeline, err := c.pipelineService.GetPipeline(id)
	if err != nil {
		// 触发成功但获取信息失败，不影响返回
		ctx.JSON(http.StatusOK, gin.H{"message": "流水线已触发", "run_id": run.RunID})
		return
	}

	// 发送通知
	go c.notifyService.SendNotification("流水线 "+pipeline.Name+" 已触发", "always")

	ctx.JSON(http.StatusOK, gin.H{"message": "流水线已触发", "run_id": run.RunID})
}

// 获取流水线日志
func (c *PipelineController) GetPipelineLogs(ctx *gin.Context) {
	runID := ctx.Param("id")
	logs, err := c.pipelineService.GetPipelineLogs(runID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "运行记录不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"logs": logs})
}

func parseUint(id uint) uint { return id }
