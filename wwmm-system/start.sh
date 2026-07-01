#!/usr/bin/env bash
# ============================================================
#  WWMM 一键启动脚本 (Linux / macOS / Git Bash)
#  用法: bash start.sh
# ============================================================
set -e
ROOT="$(cd "$(dirname "$0")" && pwd)"

MYSQL_BIN="${MYSQL_BIN:-/usr/bin/mysql}"
SERVER_PORT="${SERVER_PORT:-8080}"
FRONTEND_PORT="${FRONTEND_PORT:-5173}"

echo "============================================================"
echo " WWMM - 基于区块链的摄影作品投票存证系统"
echo " 一键启动脚本 (Linux/macOS)"
echo "============================================================"

echo ""
echo "[1/5] 检查 MySQL 连接..."
if ! "$MYSQL_BIN" -h localhost -P 3306 -u root -p123456 -e "SELECT VERSION();" >/dev/null 2>&1; then
    echo "  [X] MySQL 连接失败，请确认 MySQL 已启动"
    exit 1
fi
echo "  [OK] MySQL 已连接"

echo ""
echo "[2/5] 初始化数据库..."
"$MYSQL_BIN" -h localhost -P 3306 -u root -p123456 < "$ROOT/sql/init.sql" 2>/dev/null || echo "  [WARN] 初始化遇到警告，继续"
echo "  [OK] 数据库已初始化"

echo ""
echo "[3/5] 编译 Go 后端..."
cd "$ROOT/backend"
[ -f wwmm-server ] && rm wwmm-server
go build -o wwmm-server .
echo "  [OK] 后端已编译"

echo ""
echo "[4/5] 启动后端服务 (端口 $SERVER_PORT)..."
nohup ./wwmm-server > wwmm-server.log 2>&1 &
BACKEND_PID=$!
echo "  [OK] 后端已启动 (PID: $BACKEND_PID)"

echo ""
echo "[5/5] 启动前端服务 (端口 $FRONTEND_PORT)..."
cd "$ROOT/frontend"
[ -d node_modules ] || npm install
nohup npm run dev > vite.log 2>&1 &
FRONTEND_PID=$!
echo "  [OK] 前端已启动 (PID: $FRONTEND_PID)"

echo ""
echo "============================================================"
echo " 全部服务已启动"
echo " 后端: http://localhost:$SERVER_PORT"
echo " 前端: http://localhost:$FRONTEND_PORT"
echo ""
echo " 演示账号:"
echo "   管理员   admin / admin123"
echo "   摄影师   photographer / photo123"
echo "   投票用户 voter / vote123"
echo ""
echo " 停止服务: kill $BACKEND_PID $FRONTEND_PID"
echo "============================================================"
