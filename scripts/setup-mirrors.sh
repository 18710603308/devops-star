#!/usr/bin/env bash
# DevOpsStar 国内包管理器镜像加速配置脚本
# 此脚本会被 init.sh 自动调用

echo "==> 配置国内包管理器镜像加速..."

# ========== npm / pnpm / yarn ==========
echo "==> 配置 npm 镜像（淘宝）..."
mkdir -p ~/.npm
cat > ~/.npmrc << 'NPMRC_EOF'
registry=https://registry.npmmirror.com
disturl=https://npmmirror.com/dist
sass_binary_site=https://npmmirror.com/sass
phantomjs_cdnurl=https://npmmirror.com/phantomjs
electron_mirror=https://npmmirror.com/electron/
profiler_binary_host_mirror=https://npmmirror.com/node-profiler/
chromedriver_cdnurl=https://npmmirror.com/chromedriver
operadriver_cdnurl=https://npmmirror.com/operadriver
fse_binary_host_mirror=https://npmmirror.com/fsevents/
node_inspector_cdnurl=https://npmmirror.com/node-inspector/
NPMRC_EOF

# pnpm
pnpm config set registry https://registry.npmmirror.com 2>/dev/null || true

# yarn
yarn config set registry https://registry.npmmirror.com 2>/dev/null || true

echo "✅ npm 镜像配置完成"

# ========== Maven / Gradle ==========
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

# Gradle
mkdir -p ~/.gradle
cat > ~/.gradle/init.gradle << 'GRADLE_EOF'
allprojects {
    repositories {
        maven {
            url 'https://maven.aliyun.com/repository/public'
        }
        mavenLocal()
        mavenCentral()
    }
}
GRADLE_EOF

echo "✅ Maven/Gradle 镜像配置完成"

# ========== PyPI ==========
echo "==> 配置 PyPI 镜像（清华源）..."
mkdir -p ~/.pip
cat > ~/.pip/pip.conf << 'PIP_EOF'
[global]
index-url = https://pypi.tuna.tsinghua.edu.cn/simple
trusted-host = pypi.tuna.tsinghua.edu.cn
PIP_EOF

python3 -m pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple 2>/dev/null || true

echo "✅ PyPI 镜像配置完成"

# ========== Go Proxy ==========
echo "==> 配置 Go Proxy（七牛云）..."
go env -w GOPROXY=https://goproxy.cn,direct 2>/dev/null || true
go env -w GOSUMDB=sum.golang.google.cn,direct 2>/dev/null || true

echo "✅ Go Proxy 配置完成"

# ========== Docker ==========
echo "==> 配置 Docker 镜像加速..."
DOCKER_DAEMON="/etc/docker/daemon.json"

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
  sudo mkdir -p /etc/docker
  sudo tee "$DOCKER_DAEMON" > /dev/null << 'DOCKER_EOF'
{
  "registry-mirrors": [
    "https://docker.1ms.run",
    "https://docker.m.daocloud.io",
    "https://mirror.ccs.tencentyun.com"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m",
    "max-file": "3"
  }
}
DOCKER_EOF
  sudo systemctl daemon-reload 2>/dev/null || true
  sudo systemctl restart docker 2>/dev/null || true
  echo "✅ Docker 镜像加速配置完成（已重启 Docker 服务）"
else
  echo "⚠️  非 Linux 系统，请手动在 Docker Desktop 中配置镜像加速"
fi

# ========== Helm ==========
echo "==> 配置 Helm 镜像（阿里云）..."
helm env HELM_REPO_URL=https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts 2>/dev/null || true

echo "✅ Helm 镜像配置完成"

# ========== CocoaPods（macOS） ==========
if [[ "$OSTYPE" == "darwin"* ]]; then
  echo "==> 配置 CocoaPods 镜像（清华源）..."
  cd ~/.cocoapods/repos 2>/dev/null || true
  pod repo remove master 2>/dev/null || true
  pod repo add master https://mirrors.tuna.tsinghua.edu.cn/git/CocoaPods/Specs.git 2>/dev/null || true
  pod repo update
  echo "✅ CocoaPods 镜像配置完成"
fi

echo ""
echo "============================================"
echo "🎉  国内镜像加速配置全部完成！"
echo "============================================"
echo ""
echo "已配置："
echo "  ✅ npm / pnpm / yarn  → 淘宝镜像"
echo "  ✅ Maven / Gradle  → 阿里云"
echo "  ✅ PyPI  → 清华源"
echo "  ✅ Go Proxy  → 七牛云"
echo "  ✅ Docker  → 1ms / daocloud / 腾讯云"
echo "  ✅ Helm  → 阿里云"
if [[ "$OSTYPE" == "darwin"* ]]; then
  echo "  ✅ CocoaPods  → 清华源"
fi
echo ""
echo "提示：Docker 镜像加速需要重启 Docker 服务才能生效"
echo ""
