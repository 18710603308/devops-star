package routes

import (
	"github.com/gin-gonic/gin"
	"devops-star/backend/config"
	"devops-star/backend/controllers"
	"devops-star/backend/middleware"
)

// RegisterRoutes 注册所有 API 路由（使用真实服务）
func RegisterRoutes(
	r *gin.Engine,
	cfg *config.Config,
	authCtrl *controllers.AuthController,
	projectCtrl *controllers.ProjectController,
	pipelineCtrl *controllers.PipelineController,
	monitorCtrl *controllers.MonitorController,
	deployCtrl *controllers.DeployController,
	registryCtrl *controllers.RegistryController,
) {
	// API v1 分组
	v1 := r.Group("/api/v1")
	{
		// ========== 认证模块（无需认证）==========
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authCtrl.Login)
			auth.POST("/register", authCtrl.Register)
			auth.GET("/me", middleware.JWTAuth(cfg), authCtrl.GetMe)
		}

		// ========== 需要认证的路由 ==========
		// 使用 JWT 认证中间件
		authGroup := v1.Group("")
		authGroup.Use(middleware.JWTAuth(cfg))
		{
			// ========== 项目管理模块 ==========
			projects := authGroup.Group("/projects")
			{
				projects.GET("", projectCtrl.ListProjects)
				projects.POST("", projectCtrl.CreateProject)
				projects.GET("/:id", projectCtrl.GetProject)
				projects.PUT("/:id", projectCtrl.UpdateProject)
				projects.DELETE("/:id", projectCtrl.DeleteProject)
				// 项目成员
				projects.GET("/:id/members", projectCtrl.GetProjectMembers)
				projects.POST("/:id/members", projectCtrl.AddProjectMember)
			}

			// ========== 流水线模块 ==========
			pipelines := authGroup.Group("/pipelines")
			{
				pipelines.GET("", pipelineCtrl.ListPipelines)
				pipelines.POST("", pipelineCtrl.CreatePipeline)
				pipelines.GET("/:id", pipelineCtrl.GetPipeline)
				pipelines.POST("/:id/trigger", pipelineCtrl.TriggerPipeline)
				pipelines.GET("/:id/logs", pipelineCtrl.GetPipelineLogs)
			}

			// ========== 制品管理模块 ==========
			artifacts := authGroup.Group("/artifacts")
			{
				artifacts.GET("/images", registryCtrl.GetImages)
				artifacts.GET("/images/:repo/tags", registryCtrl.GetTags)
				artifacts.DELETE("/images/:repo/:tag", registryCtrl.DeleteImage)
				artifacts.POST("/images/:repo/:tag/scan", registryCtrl.ScanImage)
				artifacts.GET("/images/:repo/:tag/scan-report", registryCtrl.GetScanReport)
				artifacts.GET("/status", registryCtrl.GetStatus)
			}

			// ========== 部署管理模块 ==========
			deploy := authGroup.Group("/deploy")
			{
				deploy.GET("/environments", deployCtrl.GetEnvironments)
				deploy.POST("/environments", deployCtrl.CreateEnvironment)
				deploy.GET("/environments/:id/health", deployCtrl.CheckHealth)
				deploy.POST("/deploy", deployCtrl.TriggerDeploy)
				deploy.GET("/history", deployCtrl.GetDeployHistory)
				deploy.POST("/:id/rollback", deployCtrl.RollbackDeploy)
				deploy.GET("/:id/logs", deployCtrl.GetDeployLogs)
			}

			// ========== 用户管理模块 ==========
			users := authGroup.Group("/users")
			{
				users.GET("", authCtrl.GetMe) // 获取当前用户
				users.POST("", authCtrl.Register) // 管理员创建用户
				users.PUT("/:id", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "用户更新成功（待实现）"})
				})
				users.DELETE("/:id", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "用户删除成功（待实现）"})
				})
			}

			// ========== 监控模块 ==========
			monitor := authGroup.Group("/monitor")
			{
				monitor.GET("/stats", monitorCtrl.GetStats)
				monitor.GET("/pipelines", monitorCtrl.GetPipelineStats)
				monitor.GET("/deployments", monitorCtrl.GetDeployStats)
			}

			// ========== 通知测试 ==========
			notify := authGroup.Group("/notify")
			{
				notify.POST("/test", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "通知发送成功（待实现）"})
				})
			}

			// ========== 系统设置 ==========
			settings := authGroup.Group("/settings")
			{
				settings.GET("/mirrors", func(c *gin.Context) {
					c.JSON(200, gin.H{
						"docker":  []string{"https://docker.1ms.run", "https://docker.m.daocloud.io"},
						"npm":    "https://registry.npmmirror.com",
						"maven":  "https://maven.aliyun.com/repository/public",
						"pypi":   "https://pypi.tuna.tsinghua.edu.cn/simple",
						"goproxy": "https://goproxy.cn,direct",
					})
				})
				settings.PUT("/mirrors", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "镜像源配置已更新（待实现）"})
				})
			}
		}
	}

	// 健康检查（无前缀，供 docker healthcheck 使用）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "time": gin.H{"now": "2026-05-27T00:41:00+08:00"}})
	})
}
