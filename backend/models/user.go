package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Email     string         `gorm:"size:100;uniqueIndex" json:"email"`
	Role      string         `gorm:"size:20;default:'user'" json:"role"`
	Active    bool           `gorm:"default:true" json:"active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Project 项目模型
type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;uniqueIndex;not null" json:"name"`
	DisplayName string         `gorm:"size:100" json:"display_name"`
	Description string         `gorm:"type:text" json:"description"`
	RepoURL     string         `gorm:"size:255" json:"repo_url"`
	RepoType    string         `gorm:"size:20;default:'gitea'" json:"repo_type"` // gitea, github, gitlab, gitee
	GiteaID     int            `json:"gitea_id,omitempty"`
	CreatedBy   uint           `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Pipeline 流水线模型
type Pipeline struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	ProjectID   uint           `json:"project_id"`
	Description string         `gorm:"type:text" json:"description"`
	ConfigYAML  string         `gorm:"type:text" json:"config_yaml"` // 流水线 YAML 配置
	Status      string         `gorm:"size:20;default:'idle'" json:"status"` // idle, running, success, failed
	LastRunID   string         `json:"last_run_id,omitempty"`
	CreatedBy   uint           `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// PipelineRun 流水线运行记录
type PipelineRun struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	RunID      string         `gorm:"size:50;uniqueIndex" json:"run_id"`
	PipelineID uint           `json:"pipeline_id"`
	Status     string         `gorm:"size:20" json:"status"` // pending, running, success, failed
	Trigger    string         `gorm:"size:50" json:"trigger"`
	Branch     string         `gorm:"size:100" json:"branch"`
	Logs       string         `gorm:"type:text" json:"logs,omitempty"`
	StartedAt  *time.Time    `json:"started_at,omitempty"`
	FinishedAt *time.Time    `json:"finished_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// Environment 部署环境模型
type Environment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null" json:"name"`
	DisplayName string         `gorm:"size:100" json:"display_name"`
	ProjectID   uint           `json:"project_id"`
	DeployType  string         `gorm:"size:20;default:'docker'" json:"deploy_type"` // docker, k8s
	Config      string         `gorm:"type:text" json:"config"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// DeployRecord 部署记录
type DeployRecord struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	EnvironmentID uint           `json:"environment_id"`
	PipelineRunID string         `json:"pipeline_run_id,omitempty"`
	Status        string         `gorm:"size:20" json:"status"` // pending, running, success, failed, rolled_back
	DeployURL     string         `gorm:"size:255" json:"deploy_url,omitempty"`
	ImageTag      string         `gorm:"size:100" json:"image_tag,omitempty"`
	Logs          string         `gorm:"type:text" json:"logs,omitempty"`
	DeployedBy    uint           `json:"deployed_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// NotificationConfig 通知配置
type NotificationConfig struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Type      string         `gorm:"size:20;not null" json:"type"` // wecom, dingtalk, feishu, email
	Webhook  string         `gorm:"size:500" json:"webhook,omitempty"`
	SMTPHost string         `gorm:"size:100" json:"smtp_host,omitempty"`
	SMTPPort int            `json:"smtp_port,omitempty"`
	SMTPUser string         `gorm:"size:100" json:"smtp_user,omitempty"`
	SMTPPass string         `gorm:"size:255" json:"-"`
	NotifyOn []string       `gorm:"type:text" json:"notify_on,omitempty"` // success, failed, always
	Enabled  bool           `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
