# WWMM - 启动指南

## 📋 目录说明

```
.vscode/
├── tasks.json       # VSCode 任务定义（推荐）
├── launch.json      # VSCode 调试配置
├── settings.json    # VSCode 工作区设置
└── extensions.json  # 推荐的 VSCode 扩展

start.bat / start.sh  # 一键启动脚本（不依赖 VSCode）
stop.bat              # 停止 Windows 服务
```

---

## 🚀 方式一：VSCode 任务（推荐）

按 `Ctrl + Shift + P`（Mac: `Cmd + Shift + P`），输入 "Tasks: Run Task"，选择以下任务之一：

| 任务名 | 作用 |
|--------|------|
| **wwmm: Full Setup** | 首次完整启动（建库+编译+装依赖+启动前后端） |
| **wwmm: Start Backend + Frontend** | 启动前后端（首次后每次启动用这个） |
| **wwmm: Build All** | 编译后端 + 构建前端 |
| init: Run init.sql | 重建数据库 |
| build: Backend (Go) | 编译 Go 后端 |
| install: Frontend (npm) | 安装前端依赖 |
| run: Backend (Go) | 仅启动后端 |
| run: Frontend (Vite dev) | 仅启动前端 |
| kill: wwmm processes | 停止所有相关进程 |

> 也可以直接编辑 `.vscode/tasks.json` 给每个任务绑定快捷键。

### 推荐操作流程

1. **首次启动**：运行 `wwmm: Full Setup`（自动建库+编译+装依赖+启动）
2. **日常启动**：运行 `wwmm: Start Backend + Frontend`
3. **停止服务**：运行 `kill: wwmm processes`

---

## 🐞 方式二：VSCode 调试（F5）

按 `F5` 打开"运行和调试"面板（左侧虫子图标），选择：

| 配置 | 说明 |
|------|------|
| **▶ Start WWMM (Backend + Frontend)** | 一键启动（推荐） |
| Debug: Backend (Go) | 调试 Go 后端（可打断点） |
| Debug: Attach to wwmm-server | 附加到已运行的后端进程 |

> 选择 "Debug: Backend (Go)" 后按 F5 即可启动并调试后端。

---

## 📜 方式三：命令行脚本

### Windows

双击或命令行执行：

```cmd
start.bat
```

启动 MySQL → 建库 → 编译后端 → 启动前后端。

停止服务：
```cmd
stop.bat
```

### Linux / macOS / Git Bash

```bash
bash start.sh
```

---

## 🔑 演示账号

| 角色 | 账号 | 密码 |
|------|------|------|
| 管理员 | admin | admin123 |
| 摄影师 | photographer | photo123 |
| 摄影师 | alice | alice123 |
| 摄影师 | bob | bob123 |
| 投票用户 | voter | vote123 |

---

## 🛠 故障排查

| 问题 | 解决方案 |
|------|----------|
| MySQL 连接失败 | 确认 MySQL 服务已启动，确认账号密码为 root/123456 |
| 端口 8080 被占用 | 关闭已运行的后端，或修改 `config/config.go` 中的端口 |
| 端口 5173 被占用 | Vite 会自动选择下一个可用端口，查看终端输出 |
| 前端 npm install 失败 | 设置镜像 `npm config set registry https://registry.npmmirror.com` |
| Go build 失败 | 确认 Go 1.21+ 已安装：`go version` |
| 数据库表已存在 | 执行 `init: Run init.sql` 重建 |
| 想完全清理 | 停服后执行 `DROP DATABASE wwmm_db;` 再 `init: Run init.sql` |

---

## 📦 访问入口

启动成功后：

- **前端 UI**：http://localhost:5173
- **后端 API**：http://localhost:8080
- **健康检查**：http://localhost:8080/health
- **链状态 API**：http://localhost:8080/api/chain/state
