# 🚀 DevOpsStar

> 国内开箱即用的 DevOps 平台 — 图形化流水线编排 · 代码仓库集成 · 制品管理 · 监控大屏

[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go)](https://go.dev)
[![Vue 3](https://img.shields.io/badge/Vue-3.4-4FC08D?logo=vue.js)](https://vuejs.org/)
[![Element Plus](https://img.shields.io/badge/Element%20Plus-2.7-409EFF?logo=element)](https://element-plus.org/)
[![Gitea](https://img.shields.io/badge/Gitea-1.22-609926?logo=git)](https://gitea.com/)
[![License](https://img.shields.io/badge/license-MIT-green)]()

---

## ✨ 功能特性

| 模块 | 功能 |
|------|------|
| 🔐 用户认证 | JWT 登录、注册、RBAC 权限控制 |
| 📦 项目管理 | 项目 CRUD、成员管理、代码仓库对接 |
| 🔗 流水线编排 | 可视化 Pipeline 编辑器（Vue Flow）、触发记录、实时日志 |
| 📦 制品管理 | Docker 镜像仓库、版本管理、保留策略 |
| 🚀 部署管理 | 多环境管理、滚动/蓝绿/金丝雀部署 |
| 📊 监控大屏 | ECharts 构建趋势、系统资源、Grafana 仪表盘 |
| 📢 通知集成 | 企业微信 / 钉钉 / 飞书 Webhook |
| 🔧 国内优化 | Docker/npm/Maven/PyPI/Go 国内镜像加速 |

---

## 🚀 快速开始

### 前置依赖

- Docker 20.10+ / Docker Compose v2.20+
- Git Bash（Windows） / Terminal（macOS / Linux）

### 一键启动

```bash
# 克隆项目
git clone https://github.com/your-org/devops-star.git
cd devops-star

# 复制环境变量模板
cp .env.example .env
# 按需修改 .env 中的密码等配置

# 运行初始化脚本（自动配置国内镜像加速、拉取镜像、初始化数据库）
chmod +x init.sh
./init.sh

# 访问平台
# 前端：    http://localhost
# 后端 API：http://localhost:8080
# Gitea：   http://localhost:3000
# Grafana： http://localhost:3001（admin / admin123）
```

### 默认账号

| 系统 | 用户名 | 密码 |
|------|--------|------|
| DevOpsStar | admin | admin123 |
| Gitea | admin | admin123 |
| Grafana | admin | admin123 |

---

## 📐 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                      用户浏览器 (Vue 3)                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌─────────┐ │
│  │ 项目看板  │  │ 流水线编排 │  │ 镜像管理  │  │ 监控大屏 │ │
│  └──────────┘  └──────────┘  └──────────┘  └─────────┘ │
└─────────────────────────────┬───────────────────────────────┘
                              │ HTTPS
┌─────────────────────────────▼───────────────────────────────┐
│                    Nginx (反向代理 + 静态资源)                │
└─────────────────────────────┬───────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────┐
│                      Go Backend API (Gin)                    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌─────────┐ │
│  │ 项目模块  │  │ 流水线模块 │  │ 制品模块  │  │ 通知模块 │ │
│  └──────────┘  └──────────┘  └──────────┘  └─────────┘ │
└──┬──────────────────┬──────────────────┬──────────────────┘
   │                  │                  │
┌──▼─────────┐  ┌────▼─────────┐  ┌────▼─────────┐
│ PostgreSQL  │  │    Redis      │  │   Gitea      │
│  (数据库)   │  │   (缓存)     │  │  (代码仓库)   │
└────────────┘  └──────────────┘  └──────────────┘
```

---

## 📂 项目结构

```
devops-star/
├── docker-compose.yml          # 一键部署编排文件
├── .env.example               # 环境变量模板
├── init.sh                    # 初始化脚本（国内镜像优化）
├── README.md                  # 项目文档
├── CHANGELOG.md              # 版本变更记录
│
├── frontend/                 # 前端项目（Vue 3 + Element Plus）
│   ├── package.json
│   ├── vite.config.ts
│   ├── Dockerfile
│   └── src/
│       ├── main.ts
│       ├── App.vue
│       ├── router/          # 路由配置
│       ├── stores/          # Pinia 状态管理
│       ├── api/             # API 请求封装
│       └── views/           # 页面组件
│
├── backend/                  # 后端项目（Go + Gin）
│   ├── go.mod
│   ├── main.go
│   ├── Dockerfile
│   ├── config/              # 配置文件
│   ├── controllers/         # 控制器
│   ├── services/            # 业务逻辑
│   ├── models/              # 数据模型
│   ├── middleware/          # 中间件
│   └── routes/              # 路由定义
│
├── gitea/                   # Gitea 配置
├── harbor/                  # Harbor 配置
├── monitoring/              # 监控配置
│   ├── prometheus.yml       # Prometheus 配置
│   ├── loki-config.yaml    # Loki 配置
│   └── grafana/            # Grafana 仪表盘
│
├── nginx/                   # Nginx 配置
│   └── conf.d/
│       └── devops-star.conf
│
└── scripts/                 # 部署脚本
    ├── setup-mirrors.sh     # 配置国内镜像源
    ├── init-db.sql          # 数据库初始化
    └── backup.sh           # 数据备份脚本
```

---

## ⚙️ 国内网络优化

### 自动注入的镜像加速

| 类型 | 镜像源 |
|------|---------|
| Docker | `docker.1ms.run`, `docker.m.daocloud.io`, `mirror.ccs.tencentyun.com` |
| npm | `registry.npmmirror.com`（淘宝镜像） |
| Maven | `maven.aliyun.com`（阿里云） |
| PyPI | `pypi.tuna.tsinghua.edu.cn`（清华源） |
| Go | `goproxy.cn`（七牛云） |

### 本地开发环境配置国内源

**前端（npm）**：项目已包含 `frontend/.npmrc`，无需额外配置，`npm install` 自动走淘宝镜像。

**后端（Go）**：
```bash
# 方式一：全局配置（推荐，一次生效）
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=off

# 方式二：项目级（仅当前 shell）
cd backend
export GOPROXY=https://goproxy.cn,direct
go mod tidy && go run main.go
```

**Docker 构建**：`docker-compose build` 时 Dockerfile 已自动配置国内源，无需额外操作。

### 通知渠道（国内主流）

- ✅ 企业微信机器人 Webhook
- ✅ 钉钉机器人 Webhook
- ✅ 飞书机器人 Webhook
- ✅ 邮件通知（SMTP）

---

## 📊 监控与日志

| 组件 | 用途 |
|------|------|
| Prometheus | 指标采集（CPU / 内存 / 磁盘 / 构建次数） |
| Grafana | 可视化仪表盘 |
| Loki | 日志聚合查询 |
| ECharts | 前端图表（构建趋势、成功率） |

访问 `http://localhost:3001` 查看 Grafana 监控大屏。

---

## 🔧 开发指南

### 前端开发

```bash
cd frontend
npm install
npm run dev
# 访问 http://localhost:5173
```

### 后端开发

```bash
cd backend
go mod tidy
go run main.go
# API 监听 :8080
```

### 重新构建镜像

```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

---

## 📦 部署指南

### 生产环境

1. 修改 `.env` 中的密码和密钥
2. 将 `GIN_MODE` 设置为 `release`
3. 配置 Nginx SSL 证书（HTTPS）
4. 运行 `./init.sh` 启动所有服务

### 仅启用核心服务（跳过 Harbor）

```bash
# 不启动 Harbor（资源占用较高）
docker-compose up -d postgres redis gitea backend frontend
```

---

## 🚧 开发状态

### ✅ 已完成

- [x] 项目脚手架搭建
- [x] 后端 API 框架（Go + Gin）
- [x] 前端框架（Vue 3 + Element Plus）
- [x] 用户认证（JWT 登录/注册）
- [x] 项目管理 CRUD
- [x] 流水线管理 CRUD
- [x] 通知服务（企业微信/钉钉/飞书 Webhook 集成）
- [x] Docker Compose 一键部署
- [x] 国内镜像加速配置
- [x] 数据库初始化脚本
- [x] Makefile 构建脚本

### 🚧 进行中

- [ ] Gitea API 集成（创建项目时自动创建仓库）
- [ ] Vue Flow 流水线编辑器完善
- [ ] 流水线实际执行（Gitea Actions 触发）
- [ ] 部署功能实际执行（Docker / K8s）
- [ ] JWT 中间件完整实现
- [ ] 单元测试覆盖

### 📋 待实现

- [ ] RBAC 权限控制
- [ ] 流水线日志实时推送（WebSocket）
- [ ] Grafana 仪表盘配置
- [ ] 制品仓库（Harbor）集成
- [ ] API 文档（Swagger）
- [ ] 移动端适配

---

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支（`git checkout -b feature/xxx`）
3. 提交更改（`git commit -m 'feat: 添加 xxx 功能'`）
4. 推送到分支（`git push origin feature/xxx`）
5. 创建 Pull Request

---

## 📄 许可证

[MIT License](LICENSE)

---

## 🙏 致谢

- [Gitea](https://gitea.com/) — 轻量级代码仓库
- [Gin](https://gin-gonic.com/) — Go Web 框架
- [Element Plus](https://element-plus.org/) — Vue 3 UI 组件库
- [Vue Flow](https://vueflow.dev/) — 可视化流水线编辑器
- [Prometheus](https://prometheus.io/) — 监控系统
- [Grafana](https://grafana.com/) — 可视化仪表盘

---

> 如有问题或建议，欢迎提交 Issue 或 Pull Request！
