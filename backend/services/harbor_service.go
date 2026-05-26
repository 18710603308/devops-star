package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"devops-star/backend/config"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HarborService Harbor API 客户端
type HarborService struct {
	Client       *http.Client
	BaseURL      string
	Username     string
	Password     string
	ProjectName string // 默认项目名
}

// NewHarborService 创建 Harbor 客户端
func NewHarborService(cfg *config.Config) *HarborService {
	return &HarborService{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		BaseURL:      cfg.HarborURL,
		Username:     cfg.HarborAdminUser,
		Password:     cfg.HarborAdminPassword,
		ProjectName: cfg.HarborProject,
	}
}

// ========== 镜像操作 ==========

// ListRepositories 列出镜像仓库
func (s *HarborService) ListRepositories(projectName string) ([]HarborRepository, error) {
	if projectName == "" {
		projectName = s.ProjectName
	}

	url := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories",
		s.BaseURL, projectName)

	resp, err := s.doRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("列出镜像仓库失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("列出镜像仓库失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var repos []HarborRepository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return repos, nil
}

// ListTags 列出镜像标签
func (s *HarborService) ListTags(projectName, repoName string) ([]HarborTag, error) {
	if projectName == "" {
		projectName = s.ProjectName
	}

	url := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts",
		s.BaseURL, projectName, repoName)

	resp, err := s.doRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("列出镜像标签失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("列出镜像标签失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var artifacts []HarborArtifact
	if err := json.NewDecoder(resp.Body).Decode(&artifacts); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 提取标签
	var tags []HarborTag
	for _, artifact := range artifacts {
		for _, tag := range artifact.Tags {
			tags = append(tags, HarborTag{
				Name:      tag.Name,
				ArtifactID: artifact.Digest,
				PushTime:  artifact.PushTime,
				Size:       artifact.Size,
			})
		}
	}

	return tags, nil
}

// DeleteArtifact 删除镜像制品
func (s *HarborService) DeleteArtifact(projectName, repoName, reference string) error {
	if projectName == "" {
		projectName = s.ProjectName
	}

	url := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts/%s",
		s.BaseURL, projectName, repoName, reference)

	resp, err := s.doRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("删除镜像制品失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("删除镜像制品失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetArtifact 获取镜像制品详情
func (s *HarborService) GetArtifact(projectName, repoName, reference string) (*HarborArtifact, error) {
	if projectName == "" {
		projectName = s.ProjectName
	}

	url := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts/%s",
		s.BaseURL, projectName, repoName, reference)

	resp, err := s.doRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("获取镜像制品详情失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("获取镜像制品详情失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var artifact HarborArtifact
	if err := json.NewDecoder(resp.Body).Decode(&artifact); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &artifact, nil
}

// ========== 扫描操作 ==========

// ScanArtifact 扫描镜像漏洞
func (s *HarborService) ScanArtifact(projectName, repoName, reference string) error {
	if projectName == "" {
		projectName = s.ProjectName
	}

	url := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts/%s/scan",
		s.BaseURL, projectName, repoName, reference)

	resp, err := s.doRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("扫描镜像失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("扫描镜像失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetScanReport 获取镜像扫描报告
func (s *HarborService) GetScanReport(projectName, repoName, reference string) (*HarborScanReport, error) {
	if projectName == "" {
		projectName = s.ProjectName
	}

	url := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts/%s/addition/vulnerability-details",
		s.BaseURL, projectName, repoName, reference)

	resp, err := s.doRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("获取扫描报告失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("获取扫描报告失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var report HarborScanReport
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &report, nil
}

// ========== 项目操作 ==========

// ListProjects 列出项目
func (s *HarborService) ListProjects() ([]HarborProject, error) {
	url := fmt.Sprintf("%s/api/v2.0/projects", s.BaseURL)

	resp, err := s.doRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("列出项目失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("列出项目失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var projects []HarborProject
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return projects, nil
}

// CreateProject 创建项目
func (s *HarborService) CreateProject(projectName string) error {
	url := fmt.Sprintf("%s/api/v2.0/projects", s.BaseURL)

	payload := map[string]interface{}{
		"project_name": projectName,
		"public":        false,
		"storage_limit": -1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	resp, err := s.doRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建项目失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("创建项目失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteProject 删除项目
func (s *HarborService) DeleteProject(projectName string) error {
	url := fmt.Sprintf("%s/api/v2.0/projects/%s", s.BaseURL, projectName)

	resp, err := s.doRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("删除项目失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("删除项目失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// ========== 辅助函数 ==========

// doRequest 执行 HTTP 请求（带基本认证）
func (s *HarborService) doRequest(method, url string, body *bytes.Buffer) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = body
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 基本认证
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.Username, s.Password)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	req.Header.Set("Content-Type", "application/json")

	return s.Client.Do(req)
}

// ========== 数据模型 ==========

// HarborProject Harbor 项目
type HarborProject struct {
	ProjectID   int    `json:"project_id"`
	Name        string `json:"name"`
	CreationTime string `json:"creation_time"`
	Public      bool   `json:"public"`
	RepoCount   int    `json:"repo_count"`
}

// HarborRepository Harbor 镜像仓库
type HarborRepository struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ProjectID    int    `json:"project_id"`
	Description  string `json:"description"`
	PullCount    int    `json:"pull_count"`
	ArtifactCount int    `json:"artifact_count"`
	CreationTime string `json:"creation_time"`
	UpdateTime   string `json:"update_time"`
}

// HarborArtifact Harbor 镜像制品
type HarborArtifact struct {
	Digest      string         `json:"digest"`
	Tags        []HarborTagItem `json:"tags"`
	PushTime    string         `json:"push_time"`
	Size         int64          `json:"size"`
	ScanStatus   string         `json:"scan_status"`
	Vulnerabilities []HarborVulnerability `json:"vulnerabilities,omitempty"`
}

// HarborTagItem Harbor 标签项
type HarborTagItem struct {
	Name        string `json:"name"`
	PushTime    string `json:"push_time"`
	PullTime    string `json:"pull_time"`
}

// HarborTag Harbor 标签
type HarborTag struct {
	Name        string `json:"name"`
	ArtifactID  string `json:"artifact_id"`
	PushTime    string `json:"push_time"`
	Size         int64  `json:"size"`
}

// HarborScanReport Harbor 扫描报告
type HarborScanReport struct {
	Vulnerabilities []HarborVulnerability `json:"vulnerabilities"`
	Summary        HarborScanSummary    `json:"summary"`
}

// HarborVulnerability Harbor 漏洞
type HarborVulnerability struct {
	ID          string `json:"id"`
	Severity    string `json:"severity"`
	Pkg          string `json:"package"`
	Version      string `json:"version"`
	FixedVersion string `json:"fixed_version"`
	Description  string `json:"description"`
}

// HarborScanSummary Harbor 扫描摘要
type HarborScanSummary struct {
	Total     int `json:"total"`
	Fixable   int `json:"fixable"`
	Severe    int `json:"severe"`
	High       int `json:"high"`
	Medium     int `json:"medium"`
	Low        int `json:"low"`
	None       int `json:"none"`
}

// ========== registry 控制器接口 ==========

// GetRegistryController 获取制品仓库控制器（简化版）
func (s *HarborService) GetRegistryController() map[string]interface{} {
	// 返回 Harbor 连接信息
	return map[string]interface{}{
		"type":        "harbor",
		"url":         s.BaseURL,
		"project":     s.ProjectName,
		"username":    s.Username,
		"connected":   true,
	}
}

// TestConnection 测试 Harbor 连接
func (s *HarborService) TestConnection() error {
	_, err := s.ListProjects()
	if err != nil {
		return fmt.Errorf("Harbor 连接失败: %v", err)
	}
	return nil
}
