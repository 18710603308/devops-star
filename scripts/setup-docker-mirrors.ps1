# 国内 Docker 镜像加速配置脚本
# 适用系统：Windows / macOS / Linux
# 使用方法：
#   1. 以管理员身份打开 PowerShell
#   2. 执行：Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
#   3. 右键此脚本 → "使用 PowerShell 运行"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "   Docker 国内镜像加速配置脚本" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# ============ 1. 检查 Docker 是否运行 ============
Write-Host "1. 检查 Docker 状态..." -ForegroundColor Green
try {
    $dockerInfo = docker info 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "   ✗ Docker 未运行，请先启动 Docker Desktop" -ForegroundColor Red
        Read-Host "按回车退出"
        exit 1
    }
    Write-Host "   ✓ Docker 正在运行" -ForegroundColor Green
} catch {
    Write-Host "   ✗ Docker 未安装或未运行" -ForegroundColor Red
    Read-Host "按回车退出"
    exit 1
}

# ============ 2. 配置 Docker 镜像加速 ============
Write-Host "2. 配置 Docker 镜像加速..." -ForegroundColor Green

$daemonPath = "$env:USERPROFILE\.docker\daemon.json"

# 创建目录
New-Item -ItemType Directory -Path (Split-Path $daemonPath) -Force | Out-Null

# 配置内容（包含所有国内镜像源）
$daemonConfig = @{
    "builder" = @{
        "gc" = @{
            "defaultKeepStorage" = "20GB"
            "enabled" = $true
        }
    }
    "experimental" = $false
    "features" = @{
        "buildkit" = $true
    }
    "registry-mirrors" = @(
        "https://1ms.run",
        "https://docker.mirrors.ustc.edu.cn",
        "https://hub-mirror.c.163.com",
        "https://mirror.baidubce.com",
        "https://mirror.ccs.tencentyun.com",
        "https://registry.cn-hangzhou.aliyuncs.com",
        "https://registry.cn-beijing.aliyuncs.com",
        "https://registry.cn-shanghai.aliyuncs.com",
        "https://registry.cn-shenzhen.aliyuncs.com",
        "https://registry.cn-chengdu.aliyuncs.com",
        "https://dockerhub.icu",
        "https://docker.awsl9527.cn",
        "https://docker.anyhub.us.kg",
        "https://docker.1panel.live",
        "https://atomhub.openatom.cn"
    )
    "insecure-registries" = @()
}

# 写入配置文件
$daemonConfig | ConvertTo-Json -Depth 10 | Set-Content -Path $daemonPath -Encoding UTF8

Write-Host "   ✓ 配置文件已写入：$daemonPath" -ForegroundColor Green
Write-Host "   ✓ 已配置 $(($daemonConfig."registry-mirrors").Count) 个国内镜像源" -ForegroundColor Green

# ============ 3. 重启 Docker 使配置生效 ============
Write-Host "3. 重启 Docker 使配置生效..." -ForegroundColor Green

# 方法 1：通过 Docker Desktop 重启（推荐）
try {
    # 关闭 Docker Desktop
    Stop-Process -Name "Docker Desktop" -Force -ErrorAction SilentlyContinue
    Stop-Process -Name "com.docker.backend" -Force -ErrorAction SilentlyContinue
    wsl --shutdown 2>$null
    Start-Sleep -Seconds 5

    # 启动 Docker Desktop
    Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe" -ErrorAction SilentlyContinue
    
    Write-Host "   ✓ Docker Desktop 正在重启..." -ForegroundColor Green
    Write-Host "   等待 Docker 引擎启动（约 30 秒）..." -ForegroundColor Yellow

    # 等待 Docker 启动
    $retries = 0
    do {
        Start-Sleep -Seconds 5
        $retries++
        try {
            $null = docker info 2>$null
            if ($LASTEXITCODE -eq 0) { break }
        } catch {}
        Write-Host "     等待中... ($retries/12)" -ForegroundColor DarkGray
    } while ($retries -lt 12)

    if ($retries -lt 12) {
        Write-Host "   ✓ Docker 引擎已启动" -ForegroundColor Green
    } else {
        Write-Host "   ⚠ Docker 引擎启动超时，请手动检查" -ForegroundColor Yellow
    }
} catch {
    Write-Host "   ⚠ 无法自动重启 Docker，请手动重启 Docker Desktop" -ForegroundColor Yellow
}

# ============ 4. 验证配置 ============
Write-Host "4. 验证 Docker 镜像加速配置..." -ForegroundColor Green

Start-Sleep -Seconds 5

try {
    $mirrors = docker info --format "{{json .RegistryMirrors}}" 2>$null | ConvertFrom-Json
    if ($mirrors -and $mirrors.Count -gt 0) {
        Write-Host "   ✓ 镜像加速已生效！" -ForegroundColor Green
        Write-Host "   已配置的镜像源：" -ForegroundColor Cyan
        $mirrors | ForEach-Object { Write-Host "     - $_" -ForegroundColor White }
    } else {
        Write-Host "   ⚠ 镜像加速可能未生效，请手动重启 Docker Desktop" -ForegroundColor Yellow
    }
} catch {
    Write-Host "   ⚠ 无法验证配置，请手动执行：docker info" -ForegroundColor Yellow
}

# ============ 5. 拉取常用镜像（可选） ============
Write-Host "5. 是否拉取常用镜像？（y/n）" -ForegroundColor Green
$answer = Read-Host "   输入 y 拉取常用镜像，输入 n 跳过"

if ($answer -eq "y") {
    Write-Host "   开始拉取常用镜像..." -ForegroundColor Cyan
    
    $images = @(
        "nginx:alpine",
        "redis:7-alpine",
        "postgres:16-alpine",
        "node:20-alpine",
        "golang:1.22-alpine",
        "alpine:latest",
        "prom/prometheus:latest",
        "grafana/grafana:latest",
        "gitea/gitea:1.22",
        "registry:2.8"
    )

    foreach ($image in $images) {
        Write-Host "   拉取：$image" -ForegroundColor White
        docker pull $image
    }

    Write-Host "   ✓ 常用镜像拉取完成" -ForegroundColor Green
}

# ============ 完成 ============
Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "   配置完成！" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "已配置的镜像源：" -ForegroundColor Yellow
Write-Host "  • 1ms.run（推荐）" -ForegroundColor White
Write-Host "  • USTC 中科大镜像" -ForegroundColor White
Write-Host "  • 网易 163 镜像" -ForegroundColor White
Write-Host "  • 百度云镜像" -ForegroundColor White
Write-Host "  • 腾讯云镜像" -ForegroundColor White
Write-Host "  • 阿里云镜像（多区域）" -ForegroundColor White
Write-Host "  • 其他第三方镜像源" -ForegroundColor White
Write-Host ""
Write-Host "验证命令：" -ForegroundColor Cyan
Write-Host "  docker info | Select-String 'Registry Mirrors'" -ForegroundColor Gray
Write-Host ""
Read-Host "按回车退出"
