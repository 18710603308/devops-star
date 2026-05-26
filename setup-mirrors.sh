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
