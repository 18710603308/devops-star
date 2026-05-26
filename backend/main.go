package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"devops-star/backend/config"
	"devops-star/backend/controllers"
	"devops-star/backend/middleware"
	"devops-star/backend/models"
	"devops-star/backend/routes"
	"devops-star/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
	Cfg    *config.Config
)

func main() {
	// 加载配置
	Cfg = config.Load()

	// 连接数据库
	initDB()

	// 连接 Redis
	initRedis()

	// 初始化服务层
	authService := services.NewAuthService(DB, Cfg)
	projectService := services.NewProjectService(DB, Cfg, nil) // GiteaService 稍后初始化
	pipelineService := services.NewPipelineService(DB, Cfg, nil) // GiteaService 稍后初始化
	notifyService := services.NewNotificationService(Cfg)
	deployService := services.NewDeployService(DB, Cfg)
	giteaService := services.NewGiteaService(Cfg)
	harborService := services.NewHarborService(Cfg)

	// 重新初始化 projectService 和 pipelineService（注入 giteaService）
	projectService = services.NewProjectService(DB, Cfg, giteaService)
	pipelineService = services.NewPipelineService(DB, Cfg, giteaService)

	// 初始化 RBAC 服务
	rbacService := services.NewRBACService(DB)

	// 初始化控制器层
	authCtrl := controllers.NewAuthController(authService)
	projectCtrl := controllers.NewProjectController(projectService)
	pipelineCtrl := controllers.NewPipelineController(pipelineService, notifyService)
	monitorCtrl := controllers.NewMonitorController(pipelineService, deployService, Cfg.PrometheusURL)
	deployCtrl := controllers.NewDeployController(deployService, notifyService)
	registryCtrl := controllers.NewRegistryController(harborService)

	// 初始化 Gin
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// RBAC 权限中间件（在 JWT 认证之后）
	// 注意：RBACMiddleware 需要在 JWTAuth 之后使用
	// 在 routes.go 中的 authGroup 已经使用了 JWTAuth
	_ = rbacService // 避免未使用错误

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})

	// 注册路由
	routes.RegisterRoutes(r, Cfg, authCtrl, projectCtrl, pipelineCtrl, monitorCtrl, deployCtrl, registryCtrl)

	// 启动服务
	addr := fmt.Sprintf(":%d", Cfg.ServerPort)
	log.Printf("🚀 DevOpsStar Backend API 启动成功: http://localhost%s", addr)
	log.Printf("📖 API 文档: http://localhost%s/swagger/doc/index.html", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

// 初始化数据库
func initDB() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		Cfg.DBHost, Cfg.DBPort, Cfg.DBUser, Cfg.DBPassword, Cfg.DBName)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("数据库实例获取失败: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据模型
	autoMigrate()

	log.Println("✅ 数据库连接成功")
}

// 自动迁移数据模型
func autoMigrate() {
	// 先尝试自动迁移，忽略 "constraint does not exist" 错误
	err := DB.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Pipeline{},
		&models.PipelineRun{},
		&models.Environment{},
		&models.DeployRecord{},
		&models.NotificationConfig{},
		// RBAC 模型
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
	)

	// 忽略 PostgreSQL "constraint does not exist" 错误（GORM AutoMigrate 已知问题）
	if err != nil {
		errStr := err.Error()
		if !strings.Contains(errStr, "does not exist") && !strings.Contains(errStr, "42704") {
			log.Fatalf("数据模型迁移失败: %v", err)
		}
		log.Printf("⚠️ 迁移警告（已忽略）: %v", err)
	}
	log.Println("✅ 数据模型迁移完成")

	// 初始化 RBAC（创建默认角色和权限）
	if err := models.InitializeRBAC(DB); err != nil {
		log.Printf("⚠️ RBAC 初始化失败: %v", err)
	} else {
		log.Println("✅ RBAC 初始化完成")
	}
}

// 初始化 Redis
func initRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Cfg.RedisHost, Cfg.RedisPort),
		Password: Cfg.RedisPassword,
		DB:       Cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := Redis.Ping(ctx).Result(); err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}

	log.Println("✅ Redis 连接成功")
}
