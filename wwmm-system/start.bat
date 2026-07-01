@echo off
chcp 65001 >nul
setlocal

echo ============================================================
echo  WWMM - 基于区块链的摄影作品投票存证系统
echo  一键启动脚本 (Windows)
echo ============================================================
echo.

set ROOT=%~dp0
set MYSQL_BIN="C:\Program Files\MySQL\MySQL Server 8.0\bin\mysql.exe"
set SERVER_PORT=8080
set FRONTEND_PORT=5173

echo [1/5] 检查 MySQL 连接...
%MYSQL_BIN% -h localhost -P 3306 -u root -p123456 -e "SELECT VERSION();" >nul 2>&1
if errorlevel 1 (
    echo   [X] MySQL 连接失败，请确认 MySQL 已启动
    pause
    exit /b 1
)
echo   [OK] MySQL 已连接

echo.
echo [2/5] 初始化数据库（重建表）...
%MYSQL_BIN% -h localhost -P 3306 -u root -p123456 < "%ROOT%sql\init.sql" 2>nul
if errorlevel 1 (
    echo   [WARN] 数据库初始化遇到问题，尝试继续
) else (
    echo   [OK] 数据库已初始化
)

echo.
echo [3/5] 编译 Go 后端...
cd /d "%ROOT%backend"
if exist wwmm-server.exe del wwmm-server.exe
go build -o wwmm-server.exe .
if errorlevel 1 (
    echo   [X] 后端编译失败
    pause
    exit /b 1
)
echo   [OK] 后端已编译: backend\wwmm-server.exe

echo.
echo [4/5] 启动后端服务 (端口 %SERVER_PORT%)...
start "WWMM-Backend" /MIN cmd /c "wwmm-server.exe"
echo   [OK] 后端已启动

echo.
echo [5/5] 启动前端服务 (端口 %FRONTEND_PORT%)...
cd /d "%ROOT%frontend"
if not exist "node_modules" (
    echo   [INFO] 首次启动，正在安装前端依赖...
    call npm install
)
start "WWMM-Frontend" /MIN cmd /c "npm run dev"
echo   [OK] 前端已启动

echo.
echo ============================================================
echo  全部服务已启动
echo  后端: http://localhost:%SERVER_PORT%
echo  前端: http://localhost:%FRONTEND_PORT%
echo.
echo  演示账号:
echo    管理员   admin / admin123
echo    摄影师   photographer / photo123
echo    投票用户 voter / vote123
echo.
echo  关闭服务: 关闭名为 WWMM-Backend / WWMM-Frontend 的窗口
echo           或执行: stop.bat
echo ============================================================
echo.
pause
