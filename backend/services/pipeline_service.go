package services

import (
	"devops-star/backend/config"
	"devops-star/backend/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ========== 流水线服务 ==========

type PipelineService struct {
	DB           *gorm.DB
	Cfg          *config.Config
	GiteaService *GiteaService
}

func NewPipelineService(db *gorm.DB, cfg *config.Config, giteaService *GiteaService) *PipelineService {
	return &PipelineService{DB: db, Cfg: cfg, GiteaService: giteaService}
}

// 创建流水线
func (s *PipelineService) CreatePipeline(name, description, configYAML string, projectID, createdBy uint) (*models.Pipeline, error) {
	pipeline := &models.Pipeline{
		Name:        name,
		Description: description,
		ConfigYAML:  configYAML,
		ProjectID:   projectID,
		Status:      "idle",
		CreatedBy:   createdBy,
	}
	if err := s.DB.Create(pipeline).Error; err != nil {
		return nil, err
	}
	return pipeline, nil
}

// 获取流水线列表
func (s *PipelineService) ListPipelines(projectID uint) ([]models.Pipeline, error) {
	var pipelines []models.Pipeline
	query := s.DB.Where("deleted_at IS NULL")
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if err := query.Find(&pipelines).Error; err != nil {
		return nil, err
	}
	return pipelines, nil
}

// 获取单个流水线
func (s *PipelineService) GetPipeline(id uint) (*models.Pipeline, error) {
	var pipeline models.Pipeline
	if err := s.DB.First(&pipeline, id).Error; err != nil {
		return nil, err
	}
	return &pipeline, nil
}

// 触发流水线（实际执行）
func (s *PipelineService) TriggerPipeline(id uint, trigger, branch string) (*models.PipelineRun, error) {
	var pipeline models.Pipeline
	if err := s.DB.First(&pipeline, id).Error; err != nil {
		return nil, err
	}

	// 更新流水线状态
	s.DB.Model(&pipeline).Update("status", "running")

	// 创建运行记录
	run := &models.PipelineRun{
		RunID:     "run-" + generateID(),
		PipelineID: id,
		Status:     "running",
		Trigger:    trigger,
		Branch:     branch,
		StartedAt:  nil, // 实际应使用 time.Now()
	}

	if err := s.DB.Create(run).Error; err != nil {
		return nil, err
	}

	// 更新最后运行 ID
	s.DB.Model(&pipeline).Update("last_run_id", run.RunID)

	// 触发 Gitea Actions（如果配置了）
	if s.GiteaService != nil {
		// 从 pipeline.ConfigYAML 解析项目信息
		// 简化版：假设项目 ID 关联了 Gitea 仓库
		var project models.Project
		if err := s.DB.First(&project, pipeline.ProjectID).Error; err == nil {
			// 触发 Gitea Actions
			repoName := project.Name
			workflowFile := "ci.yml" // 默认 workflow 文件名
			ref := branch
			inputs := map[string]string{
				"pipeline_id":   fmt.Sprintf("%d", id),
				"pipeline_run_id": run.RunID,
			}

			go func() {
				if err := s.GiteaService.TriggerWorkflow(repoName, workflowFile, ref, inputs); err != nil {
					fmt.Printf("⚠️ 触发 Gitea Actions 失败: %v\n", err)
					// 更新运行记录状态为失败
					s.DB.Model(run).Update("status", "failed")
				} else {
					fmt.Printf("✅ Gitea Actions 触发成功: %s\n", run.RunID)
				}
			}()
		}
	}

	return run, nil
}

// 获取流水线运行日志
func (s *PipelineService) GetPipelineLogs(runID string) (string, error) {
	var run models.PipelineRun
	if err := s.DB.Where("run_id = ?", runID).First(&run).Error; err != nil {
		return "", err
	}
	return run.Logs, nil
}

// 获取监控统计
func (s *PipelineService) GetPipelineStats() (map[string]interface{}, error) {
	var total, success, failed, running int64

	s.DB.Model(&models.PipelineRun{}).Count(&total)
	s.DB.Model(&models.PipelineRun{}).Where("status = ?", "success").Count(&success)
	s.DB.Model(&models.PipelineRun{}).Where("status = ?", "failed").Count(&failed)
	s.DB.Model(&models.PipelineRun{}).Where("status = ?", "running").Count(&running)

	successRate := 0.0
	if total > 0 {
		successRate = float64(success) / float64(total) * 100
	}

	return map[string]interface{}{
		"total":       total,
		"success":     success,
		"failed":      failed,
		"running":     running,
		"successRate": successRate,
	}, nil
}

// 更新流水线运行状态（由 Webhook 回调调用）
func (s *PipelineService) UpdatePipelineRunStatus(runID, status, logs string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if logs != "" {
		updates["logs"] = logs
	}
	if status == "success" || status == "failed" {
		// 实际应使用 time.Now()
		// updates["finished_at"] = time.Now()
	}

	return s.DB.Model(&models.PipelineRun{}).Where("run_id = ?", runID).Updates(updates).Error
}

// 生成简单 ID
func generateID() string {
	return "001" // 实际应使用 uuid 或 snowflake
}

// ========== 监控服务 ==========

type MonitorService struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewMonitorService(db *gorm.DB, cfg *config.Config) *MonitorService {
	return &MonitorService{DB: db, Cfg: cfg}
}

// 获取监控统计
func (s *MonitorService) GetStats() (map[string]interface{}, error) {
	var totalProjects, totalPipelines, totalRuns int64

	s.DB.Model(&models.Project{}).Where("deleted_at IS NULL").Count(&totalProjects)
	s.DB.Model(&models.Pipeline{}).Where("deleted_at IS NULL").Count(&totalPipelines)
	s.DB.Model(&models.PipelineRun{}).Count(&totalRuns)

	return map[string]interface{}{
		"totalProjects":  totalProjects,
		"totalPipelines": totalPipelines,
		"totalRuns":      totalRuns,
	}, nil
}

// 获取部署统计
func (s *MonitorService) GetDeployStats() (map[string]interface{}, error) {
	var totalEnvs, totalDeploys, successDeploys int64

	s.DB.Model(&models.Environment{}).Where("deleted_at IS NULL").Count(&totalEnvs)
	s.DB.Model(&models.DeployRecord{}).Count(&totalDeploys)
	s.DB.Model(&models.DeployRecord{}).Where("status = ?", "success").Count(&successDeploys)

	return map[string]interface{}{
		"totalEnvs":        totalEnvs,
		"totalDeploys":     totalDeploys,
		"successDeploys":  successDeploys,
	}, nil
}

// 获取构建列表（简化版）
func (s *MonitorService) GetBuildList() ([]map[string]interface{}, error) {
	var runs []models.PipelineRun
	if err := s.DB.Order("created_at DESC").Limit(10).Find(&runs).Error; err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, run := range runs {
		result = append(result, map[string]interface{}{
			"run_id":    run.RunID,
			"pipeline_id": run.PipelineID,
			"status":      run.Status,
			"trigger":     run.Trigger,
			"branch":      run.Branch,
			"started_at":  run.StartedAt,
			"finished_at": run.FinishedAt,
		})
	}

	return result, nil
}

// ========== 通知配置服务 ==========

type NotificationConfigService struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewNotificationConfigService(db *gorm.DB, cfg *config.Config) *NotificationConfigService {
	return &NotificationConfigService{DB: db, Cfg: cfg}
}

// 获取通知配置列表
func (s *NotificationConfigService) GetConfigs() ([]models.NotificationConfig, error) {
	var configs []models.NotificationConfig
	if err := s.DB.Where("deleted_at IS NULL").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// 创建通知配置
func (s *NotificationConfigService) CreateConfig(name, notifyType, webhook, smtpHost string, smtpPort int, smtpUser, smtpPass string) (*models.NotificationConfig, error) {
	config := &models.NotificationConfig{
		Name:      name,
		Type:      notifyType,
		Webhook:   webhook,
		SMTPHost:  smtpHost,
		SMTPPort:  smtpPort,
		SMTPUser:  smtpUser,
		SMTPPass:  smtpPass,
		NotifyOn:  []string{"always"},
		Enabled:   true,
	}

	if err := s.DB.Create(config).Error; err != nil {
		return nil, err
	}
	return config, nil
}

// 更新通知配置
func (s *NotificationConfigService) UpdateConfig(id uint, updates map[string]interface{}) (*models.NotificationConfig, error) {
	var config models.NotificationConfig
	if err := s.DB.First(&config, id).Error; err != nil {
		return nil, err
	}

	if err := s.DB.Model(&config).Updates(updates).Error; err != nil {
		return nil, err
	}

	s.DB.First(&config, id)
	return &config, nil
}

// 删除通知配置
func (s *NotificationConfigService) DeleteConfig(id uint) error {
	return s.DB.Delete(&models.NotificationConfig{}, id).Error
}

// 测试通知配置
func (s *NotificationConfigService) TestConfig(id uint) error {
	var config models.NotificationConfig
	if err := s.DB.First(&config, id).Error; err != nil {
		return err
	}

	// 调用通知服务发送测试消息
	notifyService := NewNotificationService(s.Cfg)
	testMsg := "✅ DevOpsStar 通知测试成功！配置正确。"

	switch config.Type {
	case "wecom":
		return notifyService.sendWeCom(testMsg)
	case "dingtalk":
		return notifyService.sendDingTalk(testMsg)
	case "feishu":
		return notifyService.sendFeishu(testMsg)
	case "email":
		return notifyService.sendEmail(testMsg)
	default:
		return errors.New("不支持的通知类型: " + config.Type)
	}
}
