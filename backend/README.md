# DevOpsStar 后端 API

完整的 DevOps 平台后端服务，使用 Go + Gin 构建。

## 🚀 快速开始

### 前提条件

- Go 1.22+
- PostgreSQL 16+
- Redis 7+

### 本地开发

1. **复制配置**
   ```bash
   cp .env.example .env
   # 编辑 .env 填写你的配置
   ```

2. **安装依赖**
   ```bash
   make deps
   # 或
   go mod download
   ```

3. **运行服务**
   ```bash
   make run
   # 或热重载
   make dev
   ```

4. **访问 API**
   - 基础 URL: `http://localhost:8080/api/v1`
   - 健康检查: `http://localhost:8080/health`

### Docker 部署

```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

## 📝 API 端点

### 认证

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/auth/me` - 获取当前用户信息

### 项目

- `GET /api/v1/projects` - 获取项目列表
- `POST /api/v1/projects` - 创建项目
- `GET /api/v1/projects/:id` - 获取项目详情
- `PUT /api/v1/projects/:id` - 更新项目
- `DELETE /api/v1/projects/:id` - 删除项目

### 流水线

- `GET /api/v1/pipelines` - 获取流水线列表
- `POST /api/v1/pipelines` - 创建流水线
- `GET /api/v1/pipelines/:id` - 获取流水线详情
- `POST /api/v1/pipelines/:id/trigger` - 触发流水线
- `GET /api/v1/pipelines/:id/logs` - 获取流水线日志

### 部署

- `GET /api/v1/deploy/environments` - 获取环境列表
- `POST /api/v1/deploy/environments` - 创建环境
- `POST /api/v1/deploy/trigger` - 触发部署
- `GET /api/v1/deploy/history` - 获取部署历史

### 监控

- `GET /api/v1/monitor/stats` - 获取监控统计
- `GET /api/v1/monitor/pipelines` - 获取流水线统计
- `GET /api/v1/monitor/deployments` - 获取部署统计

## 🔧 配置

在 `.env` 文件中配置：

```env
# 服务器
SERVER_PORT=8080
GIN_MODE=debug

# 数据库
DB_HOST=localhost
DB_PORT=5432
DB_USER=devops
DB_PASSWORD=devops123
DB_NAME=devops_star

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-super-secret-key-change-in-production

# Gitea
GITEA_URL=http://localhost:3000
GITEA_ADMIN_USER=admin
GITEA_ADMIN_PASSWORD=admin123

# 通知配置
WECOM_WEBHOOK=https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx
DINGTALK_WEBHOOK=https://oapi.dingtalk.com/robot/send?access_token=xxx
FEISHU_WEBHOOK=https://open.feishu.cn/open-apis/bot/v2/hook/xxx

# 国内镜像加速
DOCKER_MIRROR=https://docker.1ms.run
NPM_MIRROR=https://registry.npmmirror.com
MAVEN_MIRROR=https://maven.aliyun.com/repository/public
PYPI_MIRROR=https://pypi.tuna.tsinghua.edu.cn/simple
GOPROXY_MIRROR=https://goproxy.cn,direct
```

## 🧪 测试

```bash
# 运行所有测试
make test

# 运行特定测试
go test -v ./services/...

# 查看测试覆盖率
go test -cover ./...
```

## 📦 项目结构

```
backend/
├── main.go              # 入口文件
├── config/              # 配置
├── models/              # 数据模型
├── services/            # 业务逻辑
├── controllers/         # HTTP 控制器
├── routes/              # 路由定义
├── middleware/           # 中间件
├── Makefile             # 构建脚本
├── go.mod               # Go 模块
└── .env.example        # 配置模板
```

## 🔐 认证

API 使用 JWT Bearer Token 认证：

```bash
# 登录获取 token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 使用 token 访问受保护的端点
curl http://localhost:8080/api/v1/projects \
  -H "Authorization: Bearer <your-token>"
```

## 🌐 通知集成

### 企业微信（WeCom）

1. 在企业微信后台创建群机器人
2. 复制 Webhook URL
3. 配置到 `.env` 的 `WECOM_WEBHOOK`

### 钉钉（DingTalk）

1. 在钉钉群中添加自定义机器人
2. 复制 Webhook URL
3. 配置到 `.env` 的 `DINGTALK_WEBHOOK`

### 飞书（Feishu）

1. 在飞书群中添加机器人
2. 复制 Webhook URL
3. 配置到 `.env` 的 `FEISHU_WEBHOOK`

## 🚀 生产部署

### 使用 Docker Compose（推荐）

```bash
cd ..
./init.sh
```

这将自动：
1. 配置 Docker 镜像加速
2. 配置国内包管理器镜像
3. 拉取所有镜像
4. 初始化数据库
5. 启动所有服务

### 手动部署

```bash
# 1. 构建二进制
make build

# 2. 配置 .env（生产配置）
vim .env

# 3. 运行数据库迁移
make migrate

# 4. 启动服务
./bin/devops-star

# 5. 使用 systemd 或 supervisor 保持运行
```

## 🐛 故障排除

### 数据库连接失败

```bash
# 检查 PostgreSQL 是否运行
pg_isready -h localhost -p 5432

# 检查数据库是否存在
psql -U devops -d devops_star -c "\dt"
```

### Redis 连接失败

```bash
# 检查 Redis 是否运行
redis-cli ping
```

### 端口已被占用

```bash
# 查看占用 8080 端口的进程
lsof -i :8080

# 修改 .env 中的 SERVER_PORT
vim .env
```

## 📚 更多信息

- [主项目 README](../README.md)
- [前端文档](../frontend/README.md)
- [API 文档](../docs/API.md)（待创建）

## 📄 许可证

MIT License - 详见 [LICENSE](../LICENSE)
