package services

import (
	"devops-star/backend/config"
	"devops-star/backend/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupPipelineTestDB 创建流水线测试数据库
func setupPipelineTestDB(t *testing.T) (*gorm.DB, *config.Config) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&models.Pipeline{}, &models.PipelineRun{})
	if err != nil {
		t.Fatalf("迁移数据库失败: %v", err)
	}

	cfg := &config.Config{
		GiteaURL:       "http://localhost:3000",
		GiteaAdminUser:  "admin",
		GiteaAdminPass: "admin123",
	}

	return db, cfg
}

// TestPipelineService_CreatePipeline 测试创建流水线
func TestPipelineService_CreatePipeline(t *testing.T) {
	db, cfg := setupPipelineTestDB(t)
	service := NewPipelineService(db, cfg, nil)

	tests := []struct {
		name        string
		pipelineName string
		projectID   uint
		description string
		configYAML  string
		wantErr     bool
	}{
		{
			name:        "有效流水线创建",
			pipelineName: "test-pipeline",
			projectID:   1,
			description: "测试流水线",
			configYAML:  "name: Test Pipeline\non: [push]",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipeline, err := service.CreatePipeline(tt.pipelineName, tt.projectID, tt.description, tt.configYAML)

			if tt.wantErr {
				assert.Error(t, err, "期望返回错误")
			} else {
				assert.NoError(t, err, "不期望返回错误")
				assert.NotNil(t, pipeline, "流水线不应为 nil")
				assert.Equal(t, tt.pipelineName, pipeline.Name, "流水线名称不匹配")
				assert.Equal(t, tt.projectID, pipeline.ProjectID, "项目 ID 不匹配")
			}
		})
	}
}

// TestPipelineService_GetPipeline 测试获取流水线
func TestPipelineService_GetPipeline(t *testing.T) {
	db, cfg := setupPipelineTestDB(t)
	service := NewPipelineService(db, cfg, nil)

	// 先创建一个流水线
	pipeline, err := service.CreatePipeline("get-test", 1, "获取测试", "name: Get Test")
	assert.NoError(t, err, "创建流水线应该成功")
	assert.NotNil(t, pipeline, "创建的流水线不应为 nil")

	// 测试获取存在的流水线
	gotPipeline, err := service.GetPipeline(pipeline.ID)
	assert.NoError(t, err, "获取流水线不应该出错")
	assert.NotNil(t, gotPipeline, "获取的流水线不应为 nil")
	assert.Equal(t, pipeline.ID, gotPipeline.ID, "流水线 ID 不匹配")
	assert.Equal(t, pipeline.Name, gotPipeline.Name, "流水线名称不匹配")

	// 测试获取不存在的流水线
	_, err = service.GetPipeline(999)
	assert.Error(t, err, "获取不存在的流水线应该返回错误")
}

// TestPipelineService_UpdatePipeline 测试更新流水线
func TestPipelineService_UpdatePipeline(t *testing.T) {
	db, cfg := setupPipelineTestDB(t)
	service := NewPipelineService(db, cfg, nil)

	// 先创建一个流水线
	pipeline, err := service.CreatePipeline("update-test", 1, "原始描述", "name: Update Test")
	assert.NoError(t, err, "创建流水线应该成功")

	// 更新流水线
	newDescription := "更新后的描述"
	newConfig := "name: Updated Pipeline\non: [pull_request]"
	updatedPipeline, err := service.UpdatePipeline(pipeline.ID, newDescription, newConfig)
	assert.NoError(t, err, "更新流水线不应该出错")
	assert.NotNil(t, updatedPipeline, "更新的流水线不应为 nil")
	assert.Equal(t, newDescription, updatedPipeline.Description, "更新后的描述不匹配")
	assert.Equal(t, newConfig, updatedPipeline.ConfigYAML, "更新后的配置不匹配")
}

// TestPipelineService_DeletePipeline 测试删除流水线
func TestPipelineService_DeletePipeline(t *testing.T) {
	db, cfg := setupPipelineTestDB(t)
	service := NewPipelineService(db, cfg, nil)

	// 先创建一个流水线
	pipeline, err := service.CreatePipeline("delete-test", 1, "删除测试", "name: Delete Test")
	assert.NoError(t, err, "创建流水线应该成功")

	// 删除流水线
	err = service.DeletePipeline(pipeline.ID)
	assert.NoError(t, err, "删除流水线不应该出错")

	// 验证流水线已被删除（软删除）
	_, err = service.GetPipeline(pipeline.ID)
	assert.Error(t, err, "删除后获取流水线应该返回错误")
}

// TestPipelineService_ListPipelines 测试列出流水线
func TestPipelineService_ListPipelines(t *testing.T) {
	db, cfg := setupPipelineTestDB(t)
	service := NewPipelineService(db, cfg, nil)

	// 创建多个流水线
	service.CreatePipeline("list-test-1", 1, "", "")
	service.CreatePipeline("list-test-2", 1, "", "")
	service.CreatePipeline("list-test-3", 2, "", "")

	// 列出所有流水线
	pipelines, err := service.ListPipelines(0)
	assert.NoError(t, err, "列出流水线不应该出错")
	assert.GreaterOrEqual(t, len(pipelines), 3, "应该至少有 3 个流水线")

	// 按项目 ID 列出流水线
	pipelines, err = service.ListPipelines(1)
	assert.NoError(t, err, "列出项目流水线不应该出错")
	assert.Equal(t, 2, len(pipelines), "项目 1 应该有 2 个流水线")
}

// TestPipelineService_TriggerPipeline 测试触发流水线
func TestPipelineService_TriggerPipeline(t *testing.T) {
	db, cfg := setupPipelineTestDB(t)
	service := NewPipelineService(db, cfg, nil)

	// 先创建一个流水线
	pipeline, err := service.CreatePipeline("trigger-test", 1, "触发测试", "name: Trigger Test")
	assert.NoError(t, err, "创建流水线应该成功")

	// 触发流水线
	run, err := service.TriggerPipeline(pipeline.ID, "admin", "main")
	assert.NoError(t, err, "触发流水线不应该出错")
	assert.NotNil(t, run, "运行记录不应为 nil")
	assert.Equal(t, pipeline.ID, run.PipelineID, "流水线 ID 不匹配")
	assert.Equal(t, "running", run.Status, "运行状态应该是 running")

	// 验证流水线状态已更新
	updatedPipeline, _ := service.GetPipeline(pipeline.ID)
	assert.Equal(t, "running", updatedPipeline.Status, "流水线状态应该更新为 running")
}

// TestPipelineService_GetPipelineStats 测试获取流水线统计
func TestPipelineService_GetPipelineStats(t *testing.T) {
	db, cfg := setupPipelineTestDB(t)
	service := NewPipelineService(db, cfg, nil)

	// 创建流水线和运行记录
	pipeline, _ := service.CreatePipeline("stats-test", 1, "统计测试", "")
	service.TriggerPipeline(pipeline.ID, "admin", "main")

	// 等待一秒以确保时间戳不同
	time.Sleep(1 * time.Second)

	pipeline2, _ := service.CreatePipeline("stats-test-2", 1, "统计测试2", "")
	service.TriggerPipeline(pipeline2.ID, "admin", "main")

	// 获取统计
	stats, err := service.GetPipelineStats()
	assert.NoError(t, err, "获取统计不应该出错")
	assert.NotNil(t, stats, "统计不应为 nil")
	assert.GreaterOrEqual(t, stats["total_runs"].(int64), int64(2), "总运行次数应该至少为 2")
}
