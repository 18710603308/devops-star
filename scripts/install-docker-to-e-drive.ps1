# Install-DockerToE-Drive.ps1
# 从 0 到 1 安装 Docker Desktop 并配置所有数据到 E 盘
# 用法：右键此文件 →「使用 PowerShell 运行」（需管理员权限）

# ============================================================
# 检查管理员权限
# ============================================================
if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Host "==========================================" -ForegroundColor Red
    Write-Host "  请以管理员身份运行此脚本！" -ForegroundColor Red
    Write-Host "  右键此文件 →「以管理员身份运行」" -ForegroundColor Yellow
    Write-Host "==========================================" -ForegroundColor Red
    Read-Host "按回车退出"
    exit 1
}

$ErrorActionPreference = "Stop"

# ============================================================
# 配置路径
# ============================================================
$E_DockerData   = "E:\docker-data"          # Docker 镜像/容器数据
$E_WSLRoot      = "E:\wsl"                 # WSL2 数据根目录
$E_WSLDocker     = "$E_WSLRoot\docker-desktop-data"
$E_Temp          = "E:\temp"
$DaemonConfigPath = "C:\ProgramData\Docker\config\daemon.json"

# ============================================================
# 第 1 步：创建 E 盘目录
# ============================================================
Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host " 第 1 步：创建 E 盘目录" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

New-Item -ItemType Directory -Path $E_DockerData -Force | Out-Null
New-Item -ItemType Directory -Path $E_WSLRoot  -Force | Out-Null
New-Item -ItemType Directory -Path $E_Temp     -Force | Out-Null
Write-Host "  ✓ E:\docker-data 已创建" -ForegroundColor Green
Write-Host "  ✓ E:\wsl 已创建"          -ForegroundColor Green

# ============================================================
# 第 2 步：检查 / 安装 Docker Desktop
# ============================================================
Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host " 第 2 步：检查 Docker Desktop" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

$DockerExe = "C:\Program Files\Docker\Docker\Docker Desktop.exe"
$NeedInstall = $false

if (Test-Path $DockerExe) {
    Write-Host "  ✓ Docker Desktop 已安装，跳过下载" -ForegroundColor Green
    $NeedInstall = $false
} else {
    Write-Host "  ✗ Docker Desktop 未安装，需要安装" -ForegroundColor Yellow
    $NeedInstall = $true
}

if ($NeedInstall) {
    Write-Host ""
    Write-Host "  请选择安装方式：" -ForegroundColor Cyan
    Write-Host "  [1] 自动下载并安装（国内可能较慢）" -ForegroundColor White
    Write-Host "  [2] 手动下载后安装（推荐国内用户）"   -ForegroundColor White
    $choice = Read-Host "  请输入 1 或 2"

    if ($choice -eq "1") {
        # 自动下载
        Write-Host "  正在从官网下载 Docker Desktop..." -ForegroundColor Cyan
        $InstallerPath = "$E_Temp\DockerDesktopInstaller.exe"
        try {
            Invoke-WebRequest -Uri "https://desktop.docker.com/win/stable/Docker%20Desktop%20Installer.exe" `
                -OutFile $InstallerPath -UseBasicParsing
            Write-Host "  ✓ 下载完成" -ForegroundColor Green
        } catch {
            Write-Host "  ✗ 下载失败，请手动下载：" -ForegroundColor Red
            Write-Host "    https://desktop.docker.com/win/stable/Docker%20Desktop%20Installer.exe" -ForegroundColor Yellow
            Write-Host "    下载后放到：$E_Temp\DockerDesktopInstaller.exe" -ForegroundColor Yellow
            Read-Host "  下载完成后按回车继续"
        }
    } else {
        # 手动下载提示
        Write-Host ""
        Write-Host "  请手动下载 Docker Desktop：" -ForegroundColor Cyan
        Write-Host "  https://desktop.docker.com/win/stable/Docker%20Desktop%20Installer.exe" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "  下载后：「Win + R」→ 输入 cmd → 回车" -ForegroundColor Cyan
        Write-Host "  在 CMD 里执行：" -ForegroundColor Cyan
        Write-Host "  ""$E_Temp\DockerDesktopInstaller.exe"" /quiet" -ForegroundColor Yellow
        Write-Host ""
        Read-Host "  安装完成后按回车继续"
    }

    # 执行静默安装
    $InstallerPath = "$E_Temp\DockerDesktopInstaller.exe"
    if (Test-Path $InstallerPath) {
        Write-Host "  正在静默安装 Docker Desktop（约 2-5 分钟）..." -ForegroundColor Cyan
        $proc = Start-Process -FilePath $InstallerPath -Args "/quiet" -Wait -PassThru
        Write-Host "  ✓ 安装完成" -ForegroundColor Green
    }
}

# ============================================================
# 第 3 步：等待用户启动 Docker Desktop
# ============================================================
Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host " 第 3 步：启动 Docker Desktop" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

$DockerRunning = Get-Process "Docker Desktop" -ErrorAction SilentlyContinue
if (-not $DockerRunning) {
    Write-Host "  正在启动 Docker Desktop..." -ForegroundColor Cyan
    Start-Process $DockerExe
    Write-Host "  ✓ 已启动，等待初始化（约 30 秒）..." -ForegroundColor Green
    Start-Sleep -Seconds 30
} else {
    Write-Host "  ✓ Docker Desktop 已在运行" -ForegroundColor Green
}

# 等待 docker 命令可用
Write-Host "  等待 Docker 引擎就绪..." -ForegroundColor Cyan
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

if ($retries -ge 12) {
    Write-Host "  ✗ Docker 引擎启动超时，请手动打开 Docker Desktop 等待启动完成" -ForegroundColor Red
    Read-Host "  完成后按回车继续"
}

Write-Host "  ✓ Docker 引擎已就绪" -ForegroundColor Green

# ============================================================
# 第 4 步：关闭 Docker，准备迁移 WSL 数据到 E 盘
# ============================================================
Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host " 第 4 步：迁移 WSL 数据到 E 盘" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

Write-Host "  关闭 Docker 和 WSL..." -ForegroundColor Cyan
Stop-Process -Name "Docker Desktop" -Force -ErrorAction SilentlyContinue
Stop-Process -Name "com.docker.backend" -Force -ErrorAction SilentlyContinue
wsl --shutdown
Start-Sleep -Seconds 5

# 检查 docker-desktop-data 分发是否存在
$wslList = wsl --list --quiet 2>$null
if ($wslList -match "docker-desktop-data") {
    Write-Host "  发现 docker-desktop-data，开始迁移到 E 盘..." -ForegroundColor Cyan

    # 导出
    Write-Host "    导出当前数据（可能需要几分钟）..." -ForegroundColor DarkGray
    wsl --export docker-desktop-data "$E_Temp\docker-desktop-data.tar"

    # 注销
    Write-Host "    注销原分发..." -ForegroundColor DarkGray
    wsl --unregister docker-desktop-data

    # 导入到 E 盘
    Write-Host "    导入到 E:\wsl\docker-desktop-data..." -ForegroundColor DarkGray
    wsl --import docker-desktop-data $E_WSL "$E_Temp\docker-desktop-data.tar" --version 2

    # 清理临时文件
    Remove-Item "$E_Temp\docker-desktop-data.tar" -Force -ErrorAction SilentlyContinue

    Write-Host "  ✓ WSL 数据已迁移到 $E_WSL" -ForegroundColor Green
} else {
    Write-Host "  - docker-desktop-data 不存在，跳过 WSL 迁移" -ForegroundColor Yellow
}

# ============================================================
# 第 5 步：配置 daemon.json（data-root + 国内镜像）
# ============================================================
Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host " 第 5 步：配置 Docker daemon" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

$DockerDataUnix = $E_DockerData.Replace("\", "/")

$DaemonConfig = @{
    "data-root"         = $DockerDataUnix
    "registry-mirrors"  = @(
        "https://1ms.run"
        "https://docker.mirrors.ustc.edu.cn"
        "https://hub-mirror.c.163.com"
    )
    "insecure-registries" = @("registry:5000")
}

New-Item -ItemType Directory -Path (Split-Path $DaemonConfigPath) -Force | Out-Null
$DaemonConfig | ConvertTo-Json | Out-File -FilePath $DaemonConfigPath -Encoding UTF8

Write-Host "  ✓ daemon.json 已写入：" -ForegroundColor Green
Write-Host "    $DaemonConfigPath" -ForegroundColor DarkGray
Write-Host "    data-root = $DockerDataUnix" -ForegroundColor Cyan
Write-Host "    已配置国内镜像加速" -ForegroundColor Cyan

# ============================================================
# 第 6 步：重启 Docker 使配置生效
# ============================================================
Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host " 第 6 步：重启 Docker" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

Write-Host "  正在重启 Docker Desktop..." -ForegroundColor Cyan
Stop-Process -Name "Docker Desktop" -Force -ErrorAction SilentlyContinue
wsl --shutdown
Start-Sleep -Seconds 3
Start-Process $DockerExe

Write-Host "  等待 Docker 引擎就绪..." -ForegroundColor Cyan
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

# ============================================================
# 第 7 步：验证安装
# ============================================================
Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host " 第 7 步：验证安装" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

$DockerVersion = docker --version 2>$null
Write-Host "  Docker 版本：$DockerVersion" -ForegroundColor Cyan

$DockerRoot = (docker info --format "{{.DockerRootDir}}" 2>$null)
Write-Host "  Docker 数据目录：$DockerRoot" -ForegroundColor Cyan

if ($DockerRoot -eq $DockerDataUnix) {
    Write-Host "  ✓ data-root 配置正确！" -ForegroundColor Green
} else {
    Write-Host "  ⚠ data-root 未生效，当前：$DockerRoot" -ForegroundColor Yellow
    Write-Host "    请手动重启 Docker Desktop 后再检查" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "============================================" -ForegroundColor Green
Write-Host "  ✅ Docker 安装配置完成！" -ForegroundColor Green
Write-Host "============================================" -ForegroundColor Green
Write-Host ""
Write-Host "  数据目录：  $E_DockerData"       -ForegroundColor Cyan
Write-Host "  WSL 数据：  $E_WSL"              -ForegroundColor Cyan
Write-Host "  daemon.json：$DaemonConfigPath"    -ForegroundColor Cyan
Write-Host ""
Write-Host "  下一步：在 Git Bash 里执行：" -ForegroundColor Yellow
Write-Host "  cd /d/workbuddy_workspace/hy/devops-star" -ForegroundColor White
Write-Host "  ./init.sh" -ForegroundColor White
Write-Host ""

Read-Host "按回车退出"
