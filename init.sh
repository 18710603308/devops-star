#!/usr/bin/env bash
# DevOpsStar 一键初始化脚本
# 支持：macOS / Linux / Windows (Git Bash)
# 功能：配置国内镜像加速、拉取镜像、初始化数据库、启动所有服务

set -e

# ========== 颜色定义 ==========
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ========== 打印函数 ==========
info()  { echo -e "${BLUE}[INFO]${NC} $1"; }
ok()    { echo -e "${GREEN}[OK]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
error(){ echo -e "${RED}[ERROR]${NC} $1"; }

# ========== 检测操作系统 ==========
detect_os() {
  if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "linux"
  elif [[ "$OSTYPE" == "darwin"* ]]; then
    echo "macos"
  elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
    echo "windows"
  else
    echo "unknown"
  fi
}

OS=$(detect_os)
info "检测到操作系统: $OS"

# ========== 检查依赖 ==========
check_dependencies() {
  info "检查依赖..."
  local missing=()

  if ! command -v docker &>/dev/null; then
    missing+=("docker")
  fi

  if ! command -v docker-compose &>/dev/null && ! docker compose version &>/dev/null 2>&1; then
    missing+=("docker-compose")
  fi

  if [ ${#missing[@]} -gt 0 ]; then
    error "缺少依赖: ${missing[*]}"
    error "请先安装: https://docs.docker.com/get-docker/"
    exit 1
  fi

  ok "依赖检查通过"
}

# ========== 配置 Docker 镜像加速 ==========
setup_docker_mirrors() {
  info "配置 Docker 镜像加速（国内源）..."

  local daemon_file=""
  if [[ "$OS" == "linux" ]]; then
    daemon_file="/etc/docker/daemon.json"
  elif [[ "$OS" == "macos" ]]; then
    warn "macOS Docker Desktop 请手动在 Preferences → Docker Engine 中添加镜像加速"
    warn "推荐添加："
    warn '  "registry-mirrors": ["https://docker.1ms.run", "https://docker.m.daocloud.io"]'
    return 0
  elif [[ "$OS" == "windows" ]]; then
    warn "Windows Docker Desktop 请手动在 Settings → Docker Engine 中添加镜像加速"
    warn "推荐添加："
    warn '  "registry-mirrors": ["https://docker.1ms.run", "https://docker.m.daocloud.io"]'
    return 0
  fi

  # Linux 自动配置
  if [ -f "$daemon_file" ]; then
    warn "已存在 $daemon_file，备份为 ${daemon_file}.bak"
    sudo cp "$daemon_file" "${daemon_file}.bak"
  fi

  sudo tee "$daemon_file" > /dev/null <<-'EOF'
{
  "registry-mirrors": [
    "https://docker.1ms.run",
    "https://docker.m.daocloud.io",
    "https://mirror.ccs.tencentyun.com",
    "https://hub-mirror.c.163.com",
    "https://mirror.baidubce.com"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m",
    "max-file": "3"
  }
}
EOF

  info "重启 Docker 服务以应用镜像加速..."
  sudo systemctl daemon-reload 2>/dev/null || true
  sudo systemctl restart docker 2>/dev/null || true

  ok "Docker 镜像加速配置完成"
}

# ========== 配置国内包管理器镜像 ==========
setup_package_mirrors() {
  info "生成国内包管理器镜像配置脚本..."

  cat > setup-mirrors.sh << 'MIRROR_EOF'
#!/usr/bin/env bash
# 国内包管理器镜像加速配置
# 在执行 CI/CD 流水线时，此脚本会被自动执行

echo "==> 配置 npm 镜像（淘宝）..."
npm config set registry https://registry.npmmirror.com 2>/dev/null || true

echo "==> 配置 Maven 镜像（阿里云）..."
mkdir -p ~/.m2
cat > ~/.m2/settings.xml << 'MAVEN_EOF'
<?xml version="1.0" encoding="UTF-8"?>
<settings xmlns="http://maven.apache.org/SETTINGS/1.0.0"
          xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
          xsi:schemaLocation="http://maven.apache.org/SETTINGS/1.0.0 http://maven.apache.org/xsd/settings-1.0.0.xsd">
  <mirrors>
    <mirror>
      <id>aliyun</id>
      <mirrorOf>central</mirrorOf>
      <name>阿里云 Maven 镜像</name>
      <url>https://maven.aliyun.com/repository/public</url>
    </mirror>
  </mirrors>
</settings>
MAVEN_EOF

echo "==> 配置 PyPI 镜像（清华源）..."
pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple 2>/dev/null || true

echo "==> 配置 Go Proxy（七牛云）..."
go env -w GOPROXY=https://goproxy.cn,direct 2>/dev/null || true

echo "==> 配置 Docker 镜像加速..."
mkdir -p /etc/docker
cat > /etc/docker/daemon.json << 'DOCKER_EOF'
{
  "registry-mirrors": [
    "https://docker.1ms.run",
    "https://docker.m.daocloud.io",
    "https://mirror.ccs.tencentyun.com",
    "https://hub-mirror.c.163.com",
    "https://mirror.baidubce.com"
  ]
}
DOCKER_EOF

echo "==> 配置 Helm 镜像（阿里云）..."
helm env HELM_REPO_URL=https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts 2>/dev/null || true

echo "==> 镜像配置完成！"
MIRROR_EOF

  chmod +x setup-mirrors.sh
  ok "包管理器镜像配置脚本已生成: setup-mirrors.sh"
}

# ========== 生成 .env 文件 ==========
setup_env() {
  if [ ! -f .env ]; then
    info "生成 .env 配置文件（请按需修改）..."
    cp .env.example .env
    ok ".env 文件已生成，可编辑修改默认密码等配置"
  else
    info ".env 文件已存在，跳过生成"
  fi
}

# ========== 拉取 Docker 镜像 ==========
pull_images() {
  info "拉取 Docker 镜像（使用国内镜像加速）..."
  info "这可能需要几分钟，请耐心等待..."

  # 读取 .env 中的配置
  if [ -f .env ]; then
    source .env
  fi

  # 拉取核心镜像
  local images=(
    "postgres:16-alpine"
    "redis:7-alpine"
    "gitea/gitea:1.22"
    "registry:2.8"
    "nginx:alpine"
    "prom/prometheus:latest"
    "grafana/grafana:latest"
    "grafana/loki:latest"
    "node:20-alpine"
    "golang:1.23-alpine"
  )

  for img in "${images[@]}"; do
    info "拉取镜像: $img"
    docker pull "$img" || warn "拉取 $img 失败，请检查网络"
  done

  ok "镜像拉取完成"
}

# ========== 初始化数据库 ==========
init_database() {
  info "初始化数据库..."

  # 等待 PostgreSQL 启动
  info "等待 PostgreSQL 启动..."
  local retries=30
  while [ $retries -gt 0 ]; do
    if docker exec devops-postgres pg_isready -U ${POSTGRES_USER:-devops} &>/dev/null; then
      ok "PostgreSQL 已启动"
      break
    fi
    sleep 2
    retries=$((retries - 1))
  done

  if [ $retries -eq 0 ]; then
    error "PostgreSQL 启动超时"
    exit 1
  fi

  ok "数据库初始化完成"
}

# ========== 启动服务 ==========
start_services() {
  info "启动所有服务..."

  # 读取 .env
  if [ -f .env ]; then
    source .env
  fi

  docker-compose up -d

  ok "所有服务已启动！"
}

# ========== 打印访问信息 ==========
print_info() {
  echo ""
  echo -e "${GREEN}============================================${NC}"
  echo -e "${GREEN}   DevOpsStar 平台启动成功！${NC}"
  echo -e "${GREEN}============================================${NC}"
  echo ""

  # 读取 .env
  if [ -f .env ]; then
    source .env
  fi

  local domain=${DOMAIN:-localhost}

  echo -e "${BLUE}访问地址：${NC}"
  echo -e "  前端平台:   ${GREEN}http://${domain}:${FRONTEND_PORT:-80}${NC}"
  echo -e "  后端 API:   ${GREEN}http://${domain}:${BACKEND_PORT:-8080}${NC}"
  echo -e "  Gitea:      ${GREEN}http://${domain}:${GITEA_PORT:-3000}${NC}"
  echo -e "  Harbor:      ${GREEN}http://${domain}:${HARBOR_PORT:-8081}${NC} (需取消 docker-compose.yml 中 harbor 服务的注释)"
  echo -e "  Grafana:    ${GREEN}http://${domain}:${GRAFANA_PORT:-3001}${NC}"
  echo -e "  Prometheus:  ${GREEN}http://${domain}:${PROMETHEUS_PORT:-9090}${NC}"
  echo ""

  echo -e "${BLUE}默认账号密码：${NC}"
  echo -e "  Gitea 管理员:   ${YELLOW}${GITEA_ADMIN_USER:-admin} / ${GITEA_ADMIN_PASSWORD:-admin123}${NC}"
  echo -e "  Grafana 管理员: ${YELLOW}admin / ${GRAFANA_ADMIN_PASSWORD:-admin123}${NC}"
  echo -e "  Harbor 管理员:  ${YELLOW}admin / ${HARBOR_ADMIN_PASSWORD:-Harbor12345}${NC} (若启用)"
  echo ""

  echo -e "${BLUE}快速命令：${NC}"
  echo -e "  查看服务状态:   ${YELLOW}docker-compose ps${NC}"
  echo -e "  查看服务日志:   ${YELLOW}docker-compose logs -f [服务名]${NC}"
  echo -e "  停止所有服务:   ${YELLOW}docker-compose down${NC}"
  echo -e "  重启所有服务:   ${YELLOW}docker-compose restart${NC}"
  echo ""

  echo -e "${BLUE}下一步：${NC}"
  echo -e "  1. 访问前端平台完成初始化配置"
  echo -e "  2. 在 Gitea 中创建代码仓库"
  echo -e "  3. 配置流水线实现 CI/CD"
  echo ""
}

# ========== 主流程 ==========
main() {
  echo ""
  echo -e "${BLUE}============================================${NC}"
  echo -e "${BLUE}   DevOpsStar 一键初始化脚本${NC}"
  echo -e "${BLUE}============================================${NC}"
  echo ""

  check_dependencies
  setup_docker_mirrors
  setup_package_mirrors
  setup_env
  pull_images

  info "启动服务..."
  start_services

  init_database
  print_info
}

main "$@"
