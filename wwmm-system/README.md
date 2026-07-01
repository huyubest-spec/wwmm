# 基于区块链的摄影作品投票存证系统

> 毕业设计 / 项目实训 - 湖南科技职业学院 - 区块链技术应用专业

## 项目简介

本项目是一个集**作品存证**、**人气投票**、**链上验证**于一体的去中心化平台。
每一张摄影作品上传时自动计算 SHA-256 哈希并上链，永久不可篡改；
每一次投票都作为交易写入区块链，公开透明可验证。

### 核心特性

- ⛓ **自研 PoW 区块链引擎**（Go 实现，含工作量证明 + Merkle 树 + 链式结构）
- 📷 **作品存证**：上传即上链，hash + 元数据永久保存
- 🗳 **可信投票**：每用户每作品 1 票，自动写入区块
- 🔍 **第三方校验**：任意人可通过 SHA-256 哈希独立验证作品归属
- 👨‍💼 **角色管理**：普通用户 / 摄影师 / 管理员
- 🔐 **安全**：密码 SHA-256 + 盐，Token 鉴权，CORS 跨域

### 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 2.7 + Element UI 2.15 + Axios + Vue Router |
| 前端构建 | Vite 4.5 + vite-plugin-vue2 |
| 后端 | Go 1.21 + Gin 1.9 + database/sql |
| 数据库 | MySQL 8.0 |
| 区块链 | 自研 PoW 引擎（SHA-256 + Merkle Tree） |
| 智能合约 | Solidity 0.4.25（设计稿） |
| 部署 | Windows 本地 / Linux 均可 |

## 目录结构

```
wwmm-system/
├── backend/                    Go 后端
│   ├── blockchain/             自研 PoW 区块链引擎
│   │   ├── merkle.go           Merkle 树实现
│   │   └── block.go            区块 + PoW 挖矿
│   ├── config/                 配置
│   ├── controller/             HTTP 控制器
│   ├── dao/                    数据访问层
│   ├── model/                  实体对象
│   ├── service/                业务逻辑
│   ├── utils/                  工具（crypto/auth/db/response）
│   ├── main.go                 入口
│   ├── seed.go                 种子用户
│   ├── seed_photos.go          测试数据生成器
│   ├── go.mod / go.sum         Go 模块
│   └── wwmm-server.exe         编译产物
├── frontend/                   Vue 2 前端
│   ├── src/
│   │   ├── views/              11 个页面组件
│   │   ├── router/             Vue Router 配置
│   │   ├── styles/             全局样式
│   │   ├── api.js              Axios 封装
│   │   ├── App.vue             根组件
│   │   └── main.js             入口
│   ├── index.html
│   ├── vite.config.js
│   └── package.json
├── contracts/                  Solidity 智能合约
│   ├── Roles.sol               角色库
│   ├── Admin.sol               管理员
│   ├── Photographer.sol        摄影师
│   ├── PhotoEvidence.sol       作品存证主合约
│   └── README.md               合约说明
├── sql/
│   └── init.sql                数据库初始化脚本
├── screenshots/                系统截图 13 张
├── docs/                       报告文档
│   ├── 项目实训文档-*.docx
│   ├── 毕业设计-*.docx
│   └── gen_*.py                文档生成脚本
└── README.md
```

## 快速开始

### 1. 数据库初始化

```bash
mysql -h localhost -P 3306 -u root -p123456 < sql/init.sql
```

### 2. 启动后端

```bash
cd backend/
go build -o wwmm-server.exe .
./wwmm-server.exe
# 监听 :8080
```

### 3. 启动前端

```bash
cd frontend/
npm install
npm run dev
# 监听 :5173，自动代理 /api 到 :8080
```

### 4. 访问

打开浏览器，输入 http://localhost:5173

## 演示账号

| 角色 | 账号 | 密码 |
|------|------|------|
| 管理员 | admin | admin123 |
| 摄影师 | photographer | photo123 |
| 摄影师 | alice | alice123 |
| 摄影师 | bob | bob123 |
| 投票用户 | voter | vote123 |

## API 概览

### 用户
- POST /api/user/register   注册
- POST /api/user/login      登录
- POST /api/user/logout     退出
- GET  /api/user/me         当前用户

### 作品
- GET  /api/photo/list          已审核作品列表
- GET  /api/photo/pending       待审核（管理员）
- GET  /api/photo/mine          我的作品
- GET  /api/photo/:id           作品详情
- POST /api/photo/upload        上传（multipart/form-data）
- POST /api/photo/:id/audit     审核（管理员）
- POST /api/photo/:id/vote      投票

### 区块链
- GET /api/chain/state          链状态
- GET /api/chain/blocks         区块列表
- GET /api/chain/block/:index   区块详情
- GET /api/chain/tx/:hash       交易详情
- GET /api/chain/verify/:hash   哈希校验
- GET /api/chain/txs            交易列表

## 区块链设计

### 数据结构

```
Block {
    Index       int        区块高度
    PrevHash    string     前一区块哈希
    MerkleRoot  string     Merkle 根
    Timestamp   int64      时间戳
    Nonce       int64      PoW 随机数
    Difficulty  int        挖矿难度
    Hash        string     本区块哈希
    TxCount     int        交易数
    Miner       string     打包者
}

Transaction {
    TxHash    string   交易哈希
    TxType    int      1-作品存证 2-投票
    Sender    string   发送方
    Payload   string   JSON 载荷
    Status    int      0-待打包 1-已打包
}
```

### PoW 挖矿

```
hash = SHA256(index|prevHash|merkleRoot|timestamp|nonce|difficulty|txCount|miner)
要求：hash 必须以 difficulty 个 0 开头
```

### Merkle 树

```
            root
           /    \
         h01    h23
        /  \   /  \
       h0   h1 h2  h3
       |   |  |   |
       T0  T1 T2  T3
```

## 智能合约（Solidity 设计稿）

完整代码见 `contracts/` 目录，合约架构：

```
Roles (库)
  └── Admin
        └── Photographer
              └── PhotoEvidence
```

## 文档

- [毕业设计文档](docs/毕业设计-基于区块链的摄影作品投票存证系统的设计与实现.docx)
- [项目实训文档](docs/项目实训文档-基于区块链的摄影作品投票存证系统的设计与实现.docx)

## 截图

- 01 首页作品广场
- 02 登录页
- 03 注册页
- 04 排行榜
- 05 区块链浏览器
- 06 区块详情
- 07 交易详情
- 08 哈希校验
- 09 审核中心
- 10 我的作品
- 11 上传作品
- 12 作品详情
- 13 首页-已登录

## 项目亮点

1. **完整自研 PoW 区块链引擎**：不仅使用 Go 调用库，**而是从零实现** Merkle 树、链式结构、PoW 挖矿、交易验证
2. **Solidity 智能合约双轨设计**：智能合约作为"设计稿"展示 FISCO BCOS 部署形态，Go 端等价实现保证系统可运行
3. **业务上链自动化**：上传图片 / 投票时自动打包上链，无需手动触发
4. **第三方独立验证**：任何人都能通过 SHA-256 哈希在链上验证作品归属
5. **完整 MVC 分层 + RESTful API**：代码可维护性高

## 答辩重点

- 区块链的不可篡改性如何实现（哈希链 + PoW）
- 工作量证明的具体执行过程（Mine 函数）
- Merkle 树如何保证交易完整性
- 智能合约与 Go 端实现如何对应
- 数据库设计与上链如何协同
- 前端 Vite 代理如何工作
- Token 鉴权流程

## 致谢

感谢汪铭杰、罗斌两位指导老师的悉心指导，
感谢区块链技术应用专业全体老师的教导。
