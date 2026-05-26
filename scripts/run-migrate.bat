@echo off
REM Docker 迁移到 E 盘 - 启动器
REM 自动请求管理员权限并运行 PowerShell 脚本

set "SCRIPT_PATH=%~dp0migrate-docker-to-e-drive.ps1"

echo ================================
echo   Docker 数据迁移到 E 盘
echo ================================
echo.

REM 检查管理员权限
net session >nul 2>&1
if %errorLevel% NEQ 0 (
    echo [错误] 需要管理员权限！
    echo 正在重新以管理员身份启动...
    powershell -Command "Start-Process -FilePath '%~f0' -Verb RunAs"
    exit /b
)

echo [1/3] 正在关闭 Docker 和 WSL...
taskkill /F /IM "Docker Desktop.exe" >nul 2>&1
taskkill /F /IM "com.docker.backend.exe" >nul 2>&1
wsl --shutdown >nul 2>&1
timeout /t 3 /nobreak >nul
echo OK 已关闭

echo.
echo [2/3] 创建 E 盘目录...
if not exist "E:\docker-data" mkdir "E:\docker-data"
if not exist "E:\wsl"        mkdir "E:\wsl"
if not exist "E:\temp"        mkdir "E:\temp"
echo OK 目录已创建

echo.
echo [3/3] 正在配置 Docker...
powershell -NoProfile -ExecutionPolicy Bypass -Command "& '%SCRIPT_PATH%'"
echo.

echo ================================
echo  迁移完成！
echo ================================
pause
