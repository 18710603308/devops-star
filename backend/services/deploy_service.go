package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"devops-star/backend/config"
	"devops-star/backend/models"
	"gorm.io/gorm"
)

// DeployService 部署服务
type DeployService struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewDeployService(db *gorm.DB, cfg *config.Config) *DeployService {
	return &DeployService{DB: db, Cfg: cfg}
}

// ========== 环境管理 ==========

// GetEnvironments 获取部署环境列表
func (s *DeployService) GetEnvironments(projectID uint) ([]models.Environment, error) {
	var envs []models.Environment
	query := s.DB.Where("deleted_at IS NULL")
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if err := query.Find(&envs).Error; err != nil {
		return nil, err
	}
	return envs, nil
}

// CreateEnvironment 创建部署环境
func (s *DeployService) CreateEnvironment(name, displayName, deployType string, projectID uint) (*models.Environment, error) {
	env := &models.Environment{
		Name:        name,
		DisplayName: displayName,
		DeployType:  deployType,
		ProjectID:   projectID,
		Status:      "active",
	}
	if err := s.DB.Create(env).Error; err != nil {
		return nil, err
	}
	return env, nil
}

// GetEnvironment 获取单个环境
func (s *DeployService) GetEnvironment(id uint) (*models.Environment, error) {
	var env models.Environment
	if err := s.DB.First(&env, id).Error; err != nil {
		return nil, err
	}
	return &env, nil
}

// UpdateEnvironmentStatus 更新环境状态
func (s *DeployService) UpdateEnvironmentStatus(id uint, status string) error {
	return s.DB.Model(&models.Environment{}).Where("id = ?", id).Update("status", status).Error
}

// ========== 部署执行 ==========

// TriggerDeploy 触发部署
func (s *DeployService) TriggerDeploy(envID uint, runID string, imageTag string, deployedBy uint) (*models.DeployRecord, error) {
	// 获取环境信息
	env, err := s.GetEnvironment(envID)
	if err != nil {
		return nil, fmt.Errorf("获取环境失败: %v", err)
	}

	// 创建部署记录
	record := &models.DeployRecord{
		EnvironmentID: envID,
		PipelineRunID: runID,
		ImageTag:      imageTag,
		Status:        "pending",
		DeployedBy:    deployedBy,
	}
	if err := s.DB.Create(record).Error; err != nil {
		return nil, err
	}

	// 异步执行部署
	go s.executeDeploy(record.ID, env, imageTag)

	return record, nil
}

// executeDeploy 执行部署（异步）
func (s *DeployService) executeDeploy(recordID uint, env *models.Environment, imageTag string) {
	// 更新状态为 running
	s.updateDeployStatus(recordID, "running", "")

	var err error
	switch env.DeployType {
	case "docker":
		err = s.deployDocker(env, imageTag, recordID)
	case "k8s", "kubernetes":
		err = s.deployKubernetes(env, imageTag, recordID)
	default:
		err = fmt.Errorf("不支持的部署类型: %s", env.DeployType)
	}

	if err != nil {
		s.updateDeployStatus(recordID, "failed", err.Error())
		log.Printf("❌ 部署失败 [记录#%d]: %v\n", recordID, err)
	} else {
		s.updateDeployStatus(recordID, "success", "")
		log.Printf("✅ 部署成功 [记录#%d]\n", recordID)
	}
}

// updateDeployStatus 更新部署状态
func (s *DeployService) updateDeployStatus(recordID uint, status string, errMsg string) {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == "running" {
		now := time.Now()
		updates["started_at"] = &now
	} else if status == "success" || status == "failed" {
		now := time.Now()
		updates["finished_at"] = &now
	}
	if errMsg != "" {
		updates["error_message"] = errMsg
	}
	s.DB.Model(&models.DeployRecord{}).Where("id = ?", recordID).Updates(updates)
}

// ========== Docker 部署 ==========

// deployDocker 使用 Docker 部署
func (s *DeployService) deployDocker(env *models.Environment, imageTag string, recordID uint) error {
	// 容器名称：使用环境名称
	containerName := fmt.Sprintf("devops-star-%s", env.Name)

	// 1. 拉取镜像
	log.Printf("📥 拉取镜像: %s\n", imageTag)
	cmd := exec.Command("docker", "pull", imageTag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("拉取镜像失败: %v\n输出: %s", err, string(output))
	}
	log.Printf("✅ 镜像拉取成功: %s\n", strings.TrimSpace(string(output)))

	// 2. 停止并删除旧容器（如果存在）
	log.Printf("🛑 停止旧容器: %s\n", containerName)
	exec.Command("docker", "stop", containerName).Run()
	exec.Command("docker", "rm", containerName).Run()

	// 3. 启动新容器
	log.Printf("🚀 启动新容器: %s\n", containerName)
	cmd = exec.Command("docker", "run", "-d",
		"--name", containerName,
		"-p", "8080:8080", // 默认端口映射，实际应从环境配置读取
		"--restart", "unless-stopped",
		imageTag,
	)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("启动容器失败: %v\n输出: %s", err, string(output))
	}
	containerID := strings.TrimSpace(string(output))
	log.Printf("✅ 容器启动成功: %s\n", containerID)

	// 4. 保存容器 ID 到部署记录
	s.DB.Model(&models.DeployRecord{}).Where("id = ?", recordID).Update("container_id", containerID)

	return nil
}

// ========== Kubernetes 部署 ==========

// deployKubernetes 使用 Kubernetes 部署
func (s *DeployService) deployKubernetes(env *models.Environment, imageTag string, recordID uint) error {
	// 部署名称：使用环境名称
	deploymentName := fmt.Sprintf("devops-star-%s", env.Name)
	namespace := "default" // 实际应从环境配置读取

	// 1. 检查 deployment 是否存在
	cmd := exec.Command("kubectl", "get", "deployment", deploymentName, "-n", namespace, "-o", "json")
	_, err := cmd.CombinedOutput()

	if err != nil {
		// Deployment 不存在，创建新的
		log.Printf("📦 创建新 Deployment: %s\n", deploymentName)
		return s.createK8sDeployment(env, deploymentName, namespace, imageTag, recordID)
	}

	// Deployment 存在，更新镜像
	log.Printf("🔄 更新 Deployment 镜像: %s -> %s\n", deploymentName, imageTag)
	return s.updateK8sDeployment(deploymentName, namespace, imageTag, recordID)
}

// createK8sDeployment 创建 K8s Deployment
func (s *DeployService) createK8sDeployment(env *models.Environment, name, namespace, imageTag string, recordID uint) error {
	// 生成简单的 deployment YAML
	deploymentYAML := fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s
  namespace: %s
spec:
  replicas: 1
  selector:
    matchLabels:
      app: %s
  template:
    metadata:
      labels:
        app: %s
    spec:
      containers:
      - name: app
        image: %s
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: %s-service
  namespace: %s
spec:
  selector:
    app: %s
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
`, name, namespace, name, name, imageTag, name, namespace, name)

	// 写入临时文件
	tmpFile := fmt.Sprintf("/tmp/deployment-%d.yaml", recordID)
	if err := exec.Command("bash", "-c", fmt.Sprintf("cat > %s << 'EOF'\n%s\nEOF", tmpFile, deploymentYAML)).Run(); err != nil {
		return fmt.Errorf("写入临时文件失败: %v", err)
	}

	// 应用配置
	cmd := exec.Command("kubectl", "apply", "-f", tmpFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("创建 Deployment 失败: %v\n输出: %s", err, string(output))
	}
	log.Printf("✅ Deployment 创建成功: %s\n", strings.TrimSpace(string(output)))

	// 等待部署完成
	return s.waitForK8sRollout(name, namespace, recordID)
}

// updateK8sDeployment 更新 K8s Deployment 镜像
func (s *DeployService) updateK8sDeployment(name, namespace, imageTag string, recordID uint) error {
	// 使用 kubectl set image 更新镜像
	cmd := exec.Command("kubectl", "set", "image", fmt.Sprintf("deployment/%s", name),
		fmt.Sprintf("app=%s", imageTag), "-n", namespace)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("更新镜像失败: %v\n输出: %s", err, string(output))
	}
	log.Printf("✅ 镜像更新成功: %s\n", strings.TrimSpace(string(output)))

	// 等待滚动更新完成
	return s.waitForK8sRollout(name, namespace, recordID)
}

// waitForK8sRollout 等待 K8s 滚动更新完成
func (s *DeployService) waitForK8sRollout(name, namespace string, recordID uint) error {
	cmd := exec.Command("kubectl", "rollout", "status", fmt.Sprintf("deployment/%s", name), "-n", namespace, "--timeout=300s")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("等待滚动更新失败: %v\n输出: %s", err, string(output))
	}
	log.Printf("✅ 滚动更新完成: %s\n", strings.TrimSpace(string(output)))
	return nil
}

// ========== 部署记录管理 ==========

// GetDeployHistory 获取部署历史
func (s *DeployService) GetDeployHistory(envID uint) ([]models.DeployRecord, error) {
	var records []models.DeployRecord
	query := s.DB.Order("created_at DESC")
	if envID > 0 {
		query = query.Where("environment_id = ?", envID)
	}
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// GetDeployRecord 获取单个部署记录
func (s *DeployService) GetDeployRecord(id uint) (*models.DeployRecord, error) {
	var record models.DeployRecord
	if err := s.DB.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// ========== 回滚功能 ==========

// RollbackDeploy 回滚部署
func (s *DeployService) RollbackDeploy(recordID uint) error {
	// 获取当前部署记录
	record, err := s.GetDeployRecord(recordID)
	if err != nil {
		return fmt.Errorf("获取部署记录失败: %v", err)
	}

	// 获取环境信息
	env, err := s.GetEnvironment(record.EnvironmentID)
	if err != nil {
		return fmt.Errorf("获取环境失败: %v", err)
	}

	// 获取上一次成功的部署记录
	var prevRecord models.DeployRecord
	if err := s.DB.Where("environment_id = ? AND status = ? AND id < ?", env.ID, "success", recordID).
		Order("id DESC").First(&prevRecord).Error; err != nil {
		return fmt.Errorf("未找到可回滚的历史版本: %v", err)
	}

	log.Printf("⏪ 开始回滚 [记录#%d] -> [记录#%d] 镜像: %s\n", recordID, prevRecord.ID, prevRecord.ImageTag)

	// 创建回滚记录
	rollbackRecord := &models.DeployRecord{
		EnvironmentID: record.EnvironmentID,
		PipelineRunID: prevRecord.PipelineRunID,
		ImageTag:      prevRecord.ImageTag,
		Status:        "rolling_back",
		DeployedBy:    record.DeployedBy,
	}
	if err := s.DB.Create(rollbackRecord).Error; err != nil {
		return err
	}

	// 执行回滚部署
	go func() {
		var err error
		switch env.DeployType {
		case "docker":
			err = s.deployDocker(env, prevRecord.ImageTag, rollbackRecord.ID)
		case "k8s", "kubernetes":
			err = s.deployKubernetes(env, prevRecord.ImageTag, rollbackRecord.ID)
		}

		if err != nil {
			s.updateDeployStatus(rollbackRecord.ID, "rollback_failed", err.Error())
			log.Printf("❌ 回滚失败 [记录#%d]: %v\n", rollbackRecord.ID, err)
		} else {
			s.updateDeployStatus(rollbackRecord.ID, "rolled_back", "")
			// 更新原记录状态
			s.DB.Model(&models.DeployRecord{}).Where("id = ?", recordID).Update("status", "rolled_back")
			log.Printf("✅ 回滚成功 [记录#%d]\n", rollbackRecord.ID)
		}
	}()

	return nil
}

// ========== 部署日志 ==========

// GetDeployLogs 获取部署日志（简化实现）
func (s *DeployService) GetDeployLogs(recordID uint) (string, error) {
	record, err := s.GetDeployRecord(recordID)
	if err != nil {
		return "", err
	}

	// 简化：返回部署记录的状态和错误信息
	logs := fmt.Sprintf("部署记录 #%d\n", recordID)
	logs += fmt.Sprintf("状态: %s\n", record.Status)
	logs += fmt.Sprintf("镜像: %s\n", record.ImageTag)
	logs += fmt.Sprintf("创建时间: %s\n", record.CreatedAt.Format("2006-01-02 15:04:05"))
	if record.StartedAt != nil {
		logs += fmt.Sprintf("开始时间: %s\n", record.StartedAt.Format("2006-01-02 15:04:05"))
	}
	if record.FinishedAt != nil {
		logs += fmt.Sprintf("完成时间: %s\n", record.FinishedAt.Format("2006-01-02 15:04:05"))
	}
	if record.ErrorMessage != "" {
		logs += fmt.Sprintf("错误信息: %s\n", record.ErrorMessage)
	}

	// 如果是 Docker 部署，尝试获取容器日志
	if record.ContainerID != "" {
		cmd := exec.Command("docker", "logs", "--tail", "100", record.ContainerID)
		containerLogs, err := cmd.CombinedOutput()
		if err == nil {
			logs += "\n--- 容器日志 ---\n"
			logs += string(containerLogs)
		}
	}

	return logs, nil
}

// ========== 获取部署统计 ==========

// GetDeployStats 获取部署统计
func (s *DeployService) GetDeployStats() (map[string]int64, error) {
	stats := map[string]int64{}

	// 总部署次数
	var total int64
	s.DB.Model(&models.DeployRecord{}).Count(&total)
	stats["total"] = total

	// 成功次数
	var success int64
	s.DB.Model(&models.DeployRecord{}).Where("status = ?", "success").Count(&success)
	stats["success"] = success

	// 失败次数
	var failed int64
	s.DB.Model(&models.DeployRecord{}).Where("status = ?", "failed").Count(&failed)
	stats["failed"] = failed

	// 滚动中次数
	var rolling int64
	s.DB.Model(&models.DeployRecord{}).Where("status IN ?", []string{"running", "rolling_back"}).Count(&rolling)
	stats["rolling"] = rolling

	return stats, nil
}

// ========== 健康检查 ==========

// CheckDeploymentHealth 检查部署健康状态（简化实现）
func (s *DeployService) CheckDeploymentHealth(envID uint) (string, error) {
	env, err := s.GetEnvironment(envID)
	if err != nil {
		return "unknown", err
	}

	switch env.DeployType {
	case "docker":
		return s.checkDockerHealth(env.Name)
	case "k8s", "kubernetes":
		return s.checkK8sHealth(env.Name)
	default:
		return "unknown", nil
	}
}

// checkDockerHealth 检查 Docker 容器健康状态
func (s *DeployService) checkDockerHealth(envName string) (string, error) {
	containerName := fmt.Sprintf("devops-star-%s", envName)
	cmd := exec.Command("docker", "inspect", "-f", "{{.State.Status}}", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "stopped", nil // 容器不存在
	}
	status := strings.TrimSpace(string(output))
	if status == "running" {
		return "healthy", nil
	}
	return status, nil
}

// checkK8sHealth 检查 K8s Deployment 健康状态
func (s *DeployService) checkK8sHealth(envName string) (string, error) {
	deploymentName := fmt.Sprintf("devops-star-%s", envName)
	namespace := "default"

	cmd := exec.Command("kubectl", "get", "deployment", deploymentName, "-n", namespace, "-o", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown", nil // Deployment 不存在
	}

	var deployment struct {
		Status struct {
			ReadyReplicas  int `json:"readyReplicas"`
			Replicas       int `json:"replicas"`
			UpdatedReplicas int `json:"updatedReplicas"`
		} `json:"status"`
	}

	if err := json.Unmarshal(output, &deployment); err != nil {
		return "unknown", err
	}

	if deployment.Status.ReadyReplicas == deployment.Status.Replicas {
		return "healthy", nil
	}
	return "degraded", nil
}

// ========== 通知辅助 ==========

// NotifyDeployStatus 发送部署状态通知
func (s *DeployService) NotifyDeployStatus(recordID uint, notifyService *NotificationService) {
	record, err := s.GetDeployRecord(recordID)
	if err != nil {
		return
	}

	message := fmt.Sprintf("部署 %s #%d: %s", record.ImageTag, recordID, record.Status)
	notifyService.SendNotificationAsync(message, "always")
}
