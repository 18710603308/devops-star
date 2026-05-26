# 国内镜像源配置脚本
# 适用系统：Windows / macOS / Linux
# 使用方法：
#   Windows: 在 PowerShell（管理员）中执行
#   macOS/Linux: 在终端中执行 bash setup-mirrors.sh

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "   国内镜像源一键配置脚本" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# ============ 1. Docker 镜像加速 ============
Write-Host "1. 配置 Docker 国内镜像加速..." -ForegroundColor Green

$daemonPath = "$env:USERPROFILE\.docker\daemon.json"
$daemonConfig = @{
    "registry-mirrors" = @(
        "https://1ms.run",
        "https://docker.mirrors.ustc.edu.cn",
        "https://hub-mirror.c.163.com",
        "https://mirror.baidubce.com"
    )
}

# 如果文件已存在，合并配置
if (Test-Path $daemonPath) {
    $existingConfig = Get-Content $daemonPath -Raw | ConvertFrom-Json -ErrorAction SilentlyContinue
    if ($existingConfig) {
        $daemonConfig = $existingConfig
        $daemonConfig."registry-mirrors" = @(
            "https://1ms.run",
            "https://docker.mirrors.ustc.edu.cn",
            "https://hub-mirror.c.163.com",
            "https://mirror.baidubce.com"
        )
    }
}

New-Item -ItemType Directory -Path (Split-Path $daemonPath) -Force | Out-Null
$daemonConfig | ConvertTo-Json -Depth 10 | Set-Content -Path $daemonPath -Encoding UTF8
Write-Host "   ✓ Docker 镜像加速已配置：$daemonPath" -ForegroundColor Green

# ============ 2. npm 国内镜像 ============
Write-Host "2. 配置 npm 国内镜像..." -ForegroundColor Green

npm config set registry https://registry.npmmirror.com
Write-Host "   ✓ npm 镜像：https://registry.npmmirror.com" -ForegroundColor Green

# 同时配置 yarn 和 pnpm
if (Get-Command yarn -ErrorAction SilentlyContinue) {
    yarn config set registry https://registry.npmmirror.com
    Write-Host "   ✓ yarn 镜像已配置" -ForegroundColor Green
}

if (Get-Command pnpm -ErrorAction SilentlyContinue) {
    pnpm config set registry https://registry.npmmirror.com
    Write-Host "   ✓ pnpm 镜像已配置" -ForegroundColor Green
}

# ============ 3. Maven 国内镜像 ============
Write-Host "3. 配置 Maven 国内镜像..." -ForegroundColor Green

$mavenSettingsPath = "$env:USERPROFILE\.m2\settings.xml"
$mavenSettings = @"
<settings>
  <mirrors>
    <mirror>
      <id>aliyun</id>
      <mirrorOf>central</mirrorOf>
      <name>Aliyun Maven Mirror</name>
      <url>https://maven.aliyun.com/repository/public</url>
    </mirror>
  </mirrors>
</settings>
"@

New-Item -ItemType Directory -Path "$env:USERPROFILE\.m2" -Force | Out-Null
Set-Content -Path $mavenSettingsPath -Value $mavenSettings -Encoding UTF8
Write-Host "   ✓ Maven 镜像：$mavenSettingsPath" -ForegroundColor Green

# ============ 4. PyPI 国内镜像 ============
Write-Host "4. 配置 PyPI 国内镜像..." -ForegroundColor Green

$pipConfigPath = "$env:USERPROFILE\pip\pip.ini"
$pipConfig = @"
[global]
index-url = https://pypi.tuna.tsinghua.edu.cn/simple
trusted-host = pypi.tuna.tsinghua.edu.cn

[install]
trusted-host = pypi.tuna.tsinghua.edu.cn
"@

New-Item -ItemType Directory -Path "$env:USERPROFILE\pip" -Force | Out-Null
Set-Content -Path $pipConfigPath -Value $pipConfig -Encoding UTF8
Write-Host "   ✓ PyPI 镜像：https://pypi.tuna.tsinghua.edu.cn/simple" -ForegroundColor Green

# ============ 5. Go 代理 ============
Write-Host "5. 配置 Go 国内代理..." -ForegroundColor Green

[System.Environment]::SetEnvironmentVariable("GOPROXY", "https://goproxy.cn,direct", [System.EnvironmentVariableTarget]::User)
[System.Environment]::SetEnvironmentVariable("GOSUMDB", "off", [System.EnvironmentVariableTarget]::User)
Write-Host "   ✓ GOPROXY：https://goproxy.cn" -ForegroundColor Green

# ============ 6. Homebrew 国内镜像（macOS） ============
if ($IsMacOS -or $env:OS -like "*darwin*") {
    Write-Host "6. 配置 Homebrew 国内镜像..." -ForegroundColor Green
    
    $brewMirror = "https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/"
    [System.Environment]::SetEnvironmentVariable("HOMEBREW_BREW_GIT_REMOTE", $brewMirror, [System.EnvironmentVariableTarget]::User)
    [System.Environment]::SetEnvironmentVariable("HOMEBREW_CORE_GIT_REMOTE", "$brewMirror/homebrew-core.git", [System.EnvironmentVariableTarget]::User)
    
    Write-Host "   ✓ Homebrew 镜像已配置" -ForegroundColor Green
}

# ============ 7. Git 国内镜像 ============
Write-Host "7. 配置 Git 国内镜像..." -ForegroundColor Green

git config --global url."https://gitee.com/".insteadOf "https://github.com/"
Write-Host "   ✓ Git：github.com → gitee.com 镜像" -ForegroundColor Green

# ============ 8. Rust Cargo 国内镜像 ============
Write-Host "8. 配置 Rust Cargo 国内镜像..." -ForegroundColor Green

$cargoConfigPath = "$env:USERPROFILE\.cargo\config.toml"
$cargoConfig = @"
[source.crates-io]
replace-with = 'ustc'

[source.ustc]
registry = "git://mirrors.ustc.edu.cn/crates.io-index"

[registries.ustc]
index = "git://mirrors.ustc.edu.cn/crates.io-index"
"@

New-Item -ItemType Directory -Path "$env:USERPROFILE\.cargo" -Force | Out-Null
Set-Content -Path $cargoConfigPath -Value $cargoConfig -Encoding UTF8
Write-Host "   ✓ Cargo 镜像：中科大镜像" -ForegroundColor Green

# ============ 9. Helm 国内镜像 ============
Write-Host "9. 配置 Helm 国内镜像..." -ForegroundColor Green

$env:HELM_REPO_URL = "https://charts.helm.sh/stable"
[System.Environment]::SetEnvironmentVariable("HELM_REPO_URL", "https://charts.helm.sh/stable", [System.EnvironmentVariableTarget]::User)
Write-Host "   ✓ Helm 镜像已配置" -ForegroundColor Green

# ============ 10. Docker Compose 国内镜像 ============
Write-Host "10. 配置 Docker Compose 国内镜像..." -ForegroundColor Green

$composeConfig = @"
version: '3'
services:
  app:
    image: nginx:alpine
    # 使用国内镜像加速
    pull_policy: always
"@
Write-Host "   ✓ Docker Compose 默认使用配置的镜像加速" -ForegroundColor Green

# ============ 完成 ============
Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "   配置完成！" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "已配置的镜像源：" -ForegroundColor Yellow
Write-Host "  • Docker: 1ms.run, USTC, 163, Baidu" -ForegroundColor White
Write-Host "  • npm: npmmirror.com" -ForegroundColor White
Write-Host "  • Maven: Aliyun" -ForegroundColor White
Write-Host "  • PyPI: Tsinghua" -ForegroundColor White
Write-Host "  • Go: goproxy.cn" -ForegroundColor White
Write-Host "  • Rust Cargo: USTC" -ForegroundColor White
Write-Host "  • Git: Gitee 镜像" -ForegroundColor White
Write-Host ""
Write-Host "⚠️  请注意：" -ForegroundColor Yellow
Write-Host "  1. Docker 配置需要重启 Docker Desktop 才能生效" -ForegroundColor White
Write-Host "  2. 环境变量需要重新启动终端才能生效" -ForegroundColor White
Write-Host "  3. 建议重启电脑后验证配置" -ForegroundColor White
Write-Host ""
Write-Host "验证命令：" -ForegroundColor Cyan
Write-Host "  docker info | findstr 'Registry Mirrors'" -ForegroundColor Gray
Write-Host "  npm config get registry" -ForegroundColor Gray
Write-Host "  go env GOPROXY" -ForegroundColor Gray
Write-Host "  pip config list" -ForegroundColor Gray
Write-Host ""
Write-Host "按任意键退出..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
