package services

import (
	"bytes"
	"devops-star/backend/config"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ========== Gitea 服务（API 集成）==========

type GiteaService struct {
	Cfg    *config.Config
	Client *http.Client
}

func NewGiteaService(cfg *config.Config) *GiteaService {
	return &GiteaService{
		Cfg:    cfg,
		Client: &http.Client{},
	}
}

// ========== Gitea API 模型 ==========

type GiteaCreateRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	AutoInit    bool   `json:"auto_init"`
}

type GiteaRepository struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	HTMLURL     string `json:"html_url"`
	CloneURL    string `json:"clone_url"`
	SSHURL      string `json:"ssh_url"`
	CreatedAt   string `json:"created_at"`
}

type GiteaUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type GiteaCreateHookRequest struct {
	Type   string                 `json:"type"` // gitea, gogs, slack, discord, dingtalk, etc.
	Config map[string]string      `json:"config"`
	Events []string              `json:"events"`
	Active bool                  `json:"active"`
}

// ========== Gitea API 方法 ==========

// 创建仓库
func (s *GiteaService) CreateRepo(name, description string, isPrivate bool) (*GiteaRepository, error) {
	url := fmt.Sprintf("%s/api/v1/user/repos", s.Cfg.GiteaURL)

	reqBody := GiteaCreateRepoRequest{
		Name:        name,
		Description: description,
		Private:     isPrivate,
		AutoInit:    true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.SetBasicAuth(s.Cfg.GiteaAdminUser, s.Cfg.GiteaAdminPassword)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("创建仓库失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var repo GiteaRepository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &repo, nil
}

// 获取当前用户（验证 Token）
func (s *GiteaService) GetCurrentUser() (*GiteaUser, error) {
	url := fmt.Sprintf("%s/api/v1/user", s.Cfg.GiteaURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.SetBasicAuth(s.Cfg.GiteaAdminUser, s.Cfg.GiteaAdminPassword)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("获取用户信息失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var user GiteaUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &user, nil
}

// 创建 Webhook（用于触发 CI/CD）
func (s *GiteaService) CreateWebhook(repoName, hookType, hookURL string, events []string) error {
	// 构建 API URL
	// 格式: /api/v1/repos/{owner}/{repo}/hooks
	owner := s.Cfg.GiteaAdminUser
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/hooks", s.Cfg.GiteaURL, owner, repoName)

	// 构建请求体
	reqBody := GiteaCreateHookRequest{
		Type:   hookType,
		Config: map[string]string{
			"url":          hookURL,
			"content_type": "json",
		},
		Events: events,
		Active: true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.SetBasicAuth(s.Cfg.GiteaAdminUser, s.Cfg.GiteaAdminPassword)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 成功状态码是 201 Created
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("创建 Webhook 失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// 删除仓库
func (s *GiteaService) DeleteRepo(repoName string) error {
	owner := s.Cfg.GiteaAdminUser
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s", s.Cfg.GiteaURL, owner, repoName)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.SetBasicAuth(s.Cfg.GiteaAdminUser, s.Cfg.GiteaAdminPassword)

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("删除仓库失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// 获取仓库信息
func (s *GiteaService) GetRepo(repoName string) (*GiteaRepository, error) {
	owner := s.Cfg.GiteaAdminUser
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s", s.Cfg.GiteaURL, owner, repoName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.SetBasicAuth(s.Cfg.GiteaAdminUser, s.Cfg.GiteaAdminPassword)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("获取仓库信息失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var repo GiteaRepository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &repo, nil
}

// 搜索仓库
func (s *GiteaService) SearchRepos(keyword string) ([]GiteaRepository, error) {
	url := fmt.Sprintf("%s/api/v1/repos/search?q=%s", s.Cfg.GiteaURL, keyword)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.SetBasicAuth(s.Cfg.GiteaAdminUser, s.Cfg.GiteaAdminPassword)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("搜索仓库失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	// Gitea 搜索 API 返回 { "data": [...] }
	var result struct {
		Data []GiteaRepository `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return result.Data, nil
}

// ========== 文件操作 ==========

// 提交文件到仓库（创建或更新）
func (s *GiteaService) CommitFile(repoName, filePath, content, commitMessage string) error {
	owner := s.Cfg.GiteaAdminUser
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/contents/%s",
		s.Cfg.GiteaURL, owner, repoName, filePath)

	// 先检查文件是否存在（获取 SHA）
	var existingSHA string
	{
		checkReq, err := http.NewRequest("GET", url, nil)
		if err == nil {
			checkReq.SetBasicAuth(s.Cfg.GiteaAdminUser, s.Cfg.GiteaAdminPassword)
			checkResp, err := s.Client.Do(checkReq)
			if err == nil && checkResp.StatusCode == http.StatusOK {
				var fileInfo struct {
					SHA string `json:"sha"`
				}
				json.NewDecoder(checkResp.Body).Decode(&fileInfo)
				existingSHA = fileInfo.SHA
				checkResp.Body.Close()
			}
			if checkResp != nil {
				checkResp.Body.Close()
			}
		}
	}

	// 构建请求体
	type commitRequest struct {
		Message string `json:"message"`
		Content string `json:"content"` // Base64 编码的内容
		SHA    string `json:"sha,omitempty"` // 更新时需要
		Branch string `json:"branch,omitempty"`
	}

	// 内容需要 Base64 编码
	encodedContent := base64.StdEncoding.EncodeToString([]byte(content))

	reqBody := commitRequest{
		Message: commitMessage,
		Content: encodedContent,
		Branch:  "main",
	}

	if existingSHA != "" {
		reqBody.SHA = existingSHA
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.SetBasicAuth(s.Cfg.GiteaAdminUser, s.Cfg.GiteaAdminPassword)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("提交文件失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// 读取仓库文件内容
func (s *GiteaService) GetFile(repoName, filePath string) (string, error) {
	owner := s.Cfg.GiteaAdminUser
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/contents/%s",
		s.Cfg.GiteaAdminUser, owner, repoName, filePath)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.SetBasicAuth(s.Cfg.GiteaAdminUser, s.Cfg.GiteaAdminPassword)

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("读取文件失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var fileInfo struct {
		Content  string `json:"content"` // Base64 编码的内容
		Encoding string `json:"encoding"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&fileInfo); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 解码 Base64 内容
	if fileInfo.Encoding == "base64" {
		content, err := base64.StdEncoding.DecodeString(fileInfo.Content)
		if err != nil {
			return "", fmt.Errorf("解码文件内容失败: %v", err)
		}
		return string(content), nil
	}

	return fileInfo.Content, nil
}

// 触发 Gitea Actions 流水线
func (s *GiteaService) TriggerWorkflow(repoName, workflowFile, ref string, inputs map[string]string) error {
	owner := s.Cfg.GiteaAdminUser
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/actions/workflows/%s/dispatches",
		s.Cfg.GiteaURL, owner, repoName, workflowFile)

	payload := map[string]interface{}{
		"ref":    ref,
		"inputs": inputs,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.SetBasicAuth(s.Cfg.GiteaAdminUser, s.Cfg.GiteaAdminPassword)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 成功状态码是 204 No Content
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("触发流水线失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// 获取 Actions 运行日志
func (s *GiteaService) GetWorkflowRunLogs(repoName string, runID int64) (string, error) {
	owner := s.Cfg.GiteaAdminUser
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/actions/runs/%d/logs",
		s.Cfg.GiteaURL, owner, repoName, runID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.SetBasicAuth(s.Cfg.GiteaAdminUser, s.Cfg.GiteaAdminPassword)

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("获取日志失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	logs, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取日志失败: %v", err)
	}

	return string(logs), nil
}

// 列出 Workflow 运行记录
func (s *GiteaService) ListWorkflowRuns(repoName, workflowFile string) ([]map[string]interface{}, error) {
	owner := s.Cfg.GiteaAdminUser
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/actions/workflows/%s/runs",
		s.Cfg.GiteaURL, owner, repoName, workflowFile)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.SetBasicAuth(s.Cfg.GiteaAdminUser, s.Cfg.GiteaAdminPassword)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("获取运行记录失败 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return result.Data, nil
}

// 测试 Gitea 连接
func (s *GiteaService) TestConnection() error {
	_, err := s.GetCurrentUser()
	if err != nil {
		return fmt.Errorf("Gitea 连接失败: %v", err)
	}
	return nil
}

// 获取仓库的克隆地址（HTTP 和 SSH）
func (s *GiteaService) GetRepoCloneURLs(repoName string) (httpURL, sshURL string, err error) {
	repo, err := s.GetRepo(repoName)
	if err != nil {
		return "", "", err
	}

	return repo.CloneURL, repo.SSHURL, nil
}

// 格式化仓库名称为 Gitea 格式（小写，替换特殊字符）
func FormatRepoName(name string) string {
	// 转为小写
	name = strings.ToLower(name)
	// 替换空格和特殊字符为破折号
	replaceChars := []string{" ", "_", "."}
	for _, char := range replaceChars {
		name = strings.ReplaceAll(name, char, "-")
	}
	return name
}
