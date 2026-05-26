# Migrate-DockerToE-Drive.ps1
# 将已安装的 Docker Desktop 数据迁移到 E 盘
# 用法：在管理员 PowerShell 中执行 .\migrate-docker-to-e-drive.ps1

# === 检查管理员权限 ===
$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    Write-Host "请以管理员身份运行此脚本！" -ForegroundColor Red
    Read-Host "按回车退出"
    exit 1
}

$ErrorActionPreference = "Stop"

# === 配置路径 ===
$EDockerData  = "E:\docker-data"
$EWSLRoot     = "E:\wsl"
$EWSLDocker    = "$EWSRoot\docker-desktop-data"
$ETemp         = "E:\temp"
$DaemonConfigPath = "C:\ProgramData\Docker\config\daemon.json"

# === 第 1 步：创建目录 ===
Write-Host ""
Write-Host "=== 第 1 步：创建 E 盘目录 ===" -ForegroundColor Cyan
New-Item -ItemType Directory -Path $EDockerData -Force | Out-Null
New-Item -ItemType Directory -Path $EWSRoot  -Force | Out-Null
New-Item -ItemType Directory -Path $ETemp     -Force | Out-Null
Write-Host "  OK  E:\docker-data 已创建" -ForegroundColor Green
Write-Host "  OK  E:\wsl 已创建"        -ForegroundColor Green

# === 第 2 步：关闭 Docker 和 WSL ===
Write-Host ""
Write-Host "=== 第 2 步：关闭 Docker / WSL ===" -ForegroundColor Cyan
Write-Host "  正在关闭 Docker Desktop..." -ForegroundColor Cyan
Stop-Process -Name "Docker Desktop" -Force -ErrorAction SilentlyContinue
Stop-Process -Name "com.docker.backend" -Force -ErrorAction SilentlyContinue
Start-Sleep -Seconds 3
Write-Host "  正在关闭 WSL..." -ForegroundColor Cyan
wsl --shutdown
Start-Sleep -Seconds 3
Write-Host "  OK  Docker 和 WSL 已关闭" -ForegroundColor Green

# === 第 3 步：迁移 WSL 数据 ===
Write-Host ""
Write-Host "=== 第 3 步：迁移 WSL 数据到 E 盘 ===" -ForegroundColor Cyan
$wslOutput = wsl --list --quiet 2>$null
$wslLines  = $wslOutput -split "`r?`n"
$found = $false
foreach ($line in $wslLines) {
    if ($line -match "docker-desktop-data") { $found = $true; break }
}
if ($found) {
    Write-Host "  发现 docker-desktop-data，开始迁移..." -ForegroundColor Cyan
    Write-Host "  [1/3] 导出当前数据（可能需要几分钟）..." -ForegroundColor Cyan
    wsl --export docker-desktop-data "$ETemp\docker-desktop-data.tar"
    Write-Host "  OK  导出完成" -ForegroundColor Green

    Write-Host "  [2/3] 注销原分发..." -ForegroundColor Cyan
    wsl --unregister docker-desktop-data
    Write-Host "  OK  注销完成" -ForegroundColor Green

    Write-Host "  [3/3] 导入到 E:\wsl\docker-desktop-data..." -ForegroundColor Cyan
    wsl --import docker-desktop-data $EWSLDocker "$ETemp\docker-desktop-data.tar" --version 2
    Write-Host "  OK  导入完成" -ForegroundColor Green

    Remove-Item "$ETemp\docker-desktop-data.tar" -Force -ErrorAction SilentlyContinue
    Write-Host "  OK  临时文件已清理" -ForegroundColor Green
    Write-Host "  SUCCESS  WSL 数据已迁移到 $EWSLDocker" -ForegroundColor Green
} else {
    Write-Host "  SKIP  docker-desktop-data 不存在，跳过 WSL 迁移" -ForegroundColor Yellow
}

# === 第 4 步：配置 daemon.json ===
Write-Host ""
Write-Host "=== 第 4 步：配置 daemon.json ===" -ForegroundColor Cyan
$DockerDataUnix = "/mnt/e/docker-data"
$config = @{
    "data-root"        = $DockerDataUnix
    "registry-mirrors" = @("https://1ms.run","https://docker.mirrors.ustc.edu.cn","https://hub-mirror.c.163.com")
    "insecure-registries" = @("registry:5000")
}
New-Item -ItemType Directory -Path (Split-Path $DaemonConfigPath) -Force | Out-Null
$config | ConvertTo-Json -Compress | Out-File -FilePath $DaemonConfigPath -Encoding UTF8
Write-Host "  OK  daemon.json 已写入：$DaemonConfigPath" -ForegroundColor Green
Write-Host "        data-root = $DockerDataUnix" -ForegroundColor Cyan
Write-Host "        已配置国内镜像加速"     -ForegroundColor Cyan

# === 第 5 步：重启 Docker ===
Write-Host ""
Write-Host "=== 第 5 步：重启 Docker Desktop ===" -ForegroundColor Cyan
$DockerExe = "C:\Program Files\Docker\Docker\Docker Desktop.exe"
Write-Host "  正在启动 Docker Desktop..." -ForegroundColor Cyan
Start-Process $DockerExe
Write-Host "  OK  Docker Desktop 已启动，等待引擎就绪（约 30 秒）..." -ForegroundColor Green
$retries = 0
do {
    Start-Sleep -Seconds 5
    $retries++
    try {
        $null = docker info 2>$null
        if ($LASTEXITCODE -eq 0) { break }
    } catch {}
    Write-Host "    等待中... ($retries/12)" -ForegroundColor DarkGray
} while ($retries -lt 12)

# === 第 6 步：验证 ===
Write-Host ""
Write-Host "=== 第 6 步：验证配置 ===" -ForegroundColor Cyan
try {
    $ver = docker --version 2>$null
    Write-Host "  Docker 版本：$ver" -ForegroundColor Cyan
    $root = (docker info --format "{{.DockerRootDir}}" 2>$null)
    Write-Host "  Docker 数据目录：$root" -ForegroundColor Cyan
    if ($root -eq $DockerDataUnix) {
        Write-Host "  OK  data-root 配置正确！" -ForegroundColor Green
    } else {
        Write-Host "  WARN  data-root 未生效，当前：$root" -ForegroundColor Yellow
    }
} catch {
    Write-Host "  WARN  无法获取 Docker 信息，请检查 Docker 是否正常运行" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "============================================" -ForegroundColor Green
Write-Host "  SUCCESS  Docker 数据迁移完成！" -ForegroundColor Green
Write-Host "============================================" -ForegroundColor Green
Write-Host ""
Write-Host "  数据目录：   $EDockerData"     -ForegroundColor Cyan
Write-Host "  WSL 数据：   $EWSLDocker"    -ForegroundColor Cyan
Write-Host "  daemon.json： $DaemonConfigPath" -ForegroundColor Cyan
Write-Host ""
Write-Host "  下一步：在 Git Bash 里执行：" -ForegroundColor Yellow
Write-Host "    cd /d/workbuddy_workspace/hy/devops-star" -ForegroundColor White
Write-Host "    ./init.sh" -ForegroundColor White
Write-Host ""

Read-Host "按回车退出"
