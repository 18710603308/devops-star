package controllers

import (
	"net/http"

	"devops-star/backend/services"

	"github.com/gin-gonic/gin"
)

type RegistryController struct {
	harborService *services.HarborService
}

func NewRegistryController(harborService *services.HarborService) *RegistryController {
	return &RegistryController{harborService: harborService}
}

// 获取镜像列表
func (c *RegistryController) GetImages(ctx *gin.Context) {
	repos, err := c.harborService.ListRepositories("")
	if err != nil {
		// Harbor 未启动或连接失败，返回空列表和提示
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []interface{}{},
			"message": "Harbor 镜像仓库未连接，请启动 harbor-core 服务",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": repos})
}

// 获取镜像标签
func (c *RegistryController) GetTags(ctx *gin.Context) {
	repoName := ctx.Param("repo")
	if repoName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "缺少仓库名称"})
		return
	}

	tags, err := c.harborService.ListTags("", repoName)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []interface{}{},
			"message": "Harbor 镜像仓库未连接",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tags})
}

// 删除镜像
func (c *RegistryController) DeleteImage(ctx *gin.Context) {
	repoName := ctx.Param("repo")
	tag := ctx.Param("tag")

	if repoName == "" || tag == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "缺少仓库名称或标签"})
		return
	}

	err := c.harborService.DeleteArtifact("", repoName, tag)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Harbor 未连接: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "镜像删除成功"})
}

// 扫描镜像
func (c *RegistryController) ScanImage(ctx *gin.Context) {
	repoName := ctx.Param("repo")
	tag := ctx.Param("tag")

	if repoName == "" || tag == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "缺少仓库名称或标签"})
		return
	}

	err := c.harborService.ScanArtifact("", repoName, tag)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Harbor 未连接: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "镜像扫描已触发"})
}

// 获取扫描报告
func (c *RegistryController) GetScanReport(ctx *gin.Context) {
	repoName := ctx.Param("repo")
	tag := ctx.Param("tag")

	if repoName == "" || tag == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "缺少仓库名称或标签"})
		return
	}

	report, err := c.harborService.GetScanReport("", repoName, tag)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Harbor 未连接: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

// 获取 Harbor 连接状态
func (c *RegistryController) GetStatus(ctx *gin.Context) {
	err := c.harborService.TestConnection()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": "disconnected", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "connected"})
}
