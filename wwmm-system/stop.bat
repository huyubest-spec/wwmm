@echo off
chcp 65001 >nul
echo 正在停止 WWMM 服务...
taskkill /F /IM wwmm-server.exe 2>nul
taskkill /F /FI "WINDOWTITLE eq WWMM-Backend*" 2>nul
taskkill /F /FI "WINDOWTITLE eq WWMM-Frontend*" 2>nul
echo.
echo 已停止所有 wwmm-server 进程
echo （如需同时关闭 Vite 调试窗口，请手动关闭 "WWMM-Frontend" 窗口）
pause
