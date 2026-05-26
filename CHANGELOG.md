# 变更日志

本文档记录 DevOpsStar 项目的所有重要变更。

## [未发布] - 2026-05-27

### ✨ 新增功能
- 新增真正的 JWT 认证（使用 golang-jwt/jwt/v5 库生成和验证 token）
- 新增 JWT 认证中间件（保护 API 端点）
- 新增 Gitea API 集成服务（`services/gitea_service.go`）
- 新增项目创建时自动在 Gitea 中创建代码仓库
- 新增真正的通知服务（企业微信/钉钉/飞书 Webhook HTTP 调用）
- 新增 Vue Flow 流水线编辑器集成（拖拽节点、配置、YAML 预览）
- 新增后端 Makefile（构建、运行、测试、代码检查）
- 新增后端 README.md 文档

### 🔧 优化
- 优化项目结构：将通知服务独立为 `services/notification_service.go`
- 优化部署服务：添加 `services/deploy_service.go` 真正实现部署逻辑
- 优化路由配置：使用 JWT 中间件保护 API（登录/注册除外）
- 优化控制器：部署控制器使用真正的服务层而非模拟数据

### 🐛 修复
- 修复 `pipeline_controller.go` 第83行 bug（pipeline 变量未定义）
- 修复 `main.go` 缺少 `models` 导入导致的编译错误
- 修复 `auth_service.go` 缺少 `fmt` 导入
- 修复 `project_service.go` 中 `NotificationService` 重复定义（移至独立文件）

### 📖 文档
- 新增后端 README.md（API 端点、本地开发、配置、通知集成）
- 更新主 README.md（添加"开发状态"部分）
- 更新 CHANGELOG.md（记录本次变更）

---

## [1.0.0] - 2026-05-27

### ✨ 首次发布
- 初始化项目结构
- 集成 Vue 3 + Element Plus 前端框架
- 集成 Go + Gin 后端框架
- 集成 Gitea（代码仓库 + CI/CD）
- 集成 PostgreSQL + Redis
- 集成 Prometheus + Grafana + Loki 监控
- 集成 Harbor 镜像仓库（可选）
- 完成国内镜像加速配置
- 完成 Docker Compose 一键部署

---

## 版本号规则

本项目遵循 [语义化版本 2.0.0](https://semver.org/lang/zh-CN/)：

- **主版本号**：不兼容的 API 修改
- **次版本号**：向下兼容的功能性新增
- **修订号**：向下兼容的问题修正

## 类型说明

| 类型 | 说明 |
|------|------|
| ✨ 新增功能 | 新功能 |
| 🔧 优化 | 现有功能的改进 |
| 🐛 修复 | Bug 修复 |
| 📖 文档 | 文档更新 |
| 🔨 重构 | 代码重构 |
| ⚡ 性能 | 性能优化 |
| 🏗 构建 | 构建系统或外部依赖变更 |
| 🔒 安全 | 安全问题修复 |
