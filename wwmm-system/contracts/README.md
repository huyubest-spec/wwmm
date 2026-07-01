# 智能合约说明 (PhotoEvidence.sol)

> 摄影作品投票存证系统的核心 Solidity 智能合约

## 1. 合约架构

```
                ┌──────────────────────┐
                │       Admin          │   平台管理员角色管理
                │  (isAdmin/addAdmin)  │
                └──────────┬───────────┘
                           │ 继承
                ┌──────────▼───────────┐
                │    Photographer      │   摄影师角色管理
                │ (register/getInfo)   │
                └──────────┬───────────┘
                           │ 继承
                ┌──────────▼───────────┐
                │    PhotoEvidence     │   核心：存证 + 投票
                │  submitPhoto/voteFor │
                └──────────────────────┘
```

合约间通过 `is` 关键字形成继承链，复用 `Roles` 库的 RBAC 能力。

## 2. 关键接口

### 2.1 摄影师注册

```solidity
function register(string realName, string phone) public
```

任何以太坊地址都可以注册成为摄影师。注册成功后，该地址可调用
`submitPhoto` 上传作品。

### 2.2 提交作品存证

```solidity
function submitPhoto(
    string title,
    string imageHash,    // 64 字符的 SHA-256 十六进制
    string description,
    string category,
    string shootLocation
) public onlyPhotographer returns (uint256 photoId)
```

**关键安全约束**：
- `imageHash` 长度必须等于 64（标准 SHA-256 十六进制字符串）
- 同一 imageHash 不可重复提交（`_hashIndex` 防止重复）
- 同一摄影师上传时计数加 1，便于统计

### 2.3 管理员审核

```solidity
function auditPhoto(uint256 photoId, bool approve, string comment)
public onlyAdmin
```

仅管理员可调用。审核通过后，作品状态变为 `Approved`，允许投票。

### 2.4 投票

```solidity
function voteFor(uint256 photoId) public
```

- 仅 `Approved` 状态的作品可被投票
- 不可给自己投票
- 同一地址对同一作品只允许 1 票
- 投票后写 `Vote` 数组，并更新 photo.voteCount

### 2.5 排行榜

```solidity
function getRanking(uint256 topN) public view returns (uint256[])
```

返回得票数排名前 N 的作品 ID 列表，按得票降序排列。

## 3. 事件

合约通过 `event` 暴露关键状态变更，前端或后端可以监听：

| 事件 | 触发时机 | 参数 |
|------|---------|------|
| `PhotographerRegistered` | 摄影师注册成功 | account, realName |
| `PhotoSubmitted` | 作品提交 | photoId, photographer, imageHash, submitTime |
| `PhotoAudited` | 审核完成 | photoId, status, comment |
| `Voted` | 成功投票 | photoId, voter, voteCount |

## 4. 部署与编译

```bash
# 安装 solc 0.4.25
npm install -g solc@0.4.25

# 编译
solcjs --bin --abi --optimize PhotoEvidence.sol -o build/

# 使用 truffle / hardhat / FISCO BCOS 部署
# FISCO BCOS 推荐版本: 2.8.0+，Solidity 0.4.25
```

## 5. 在本项目中的实际作用

> **说明**：由于本地 Windows 环境部署 FISCO BCOS 4 节点联盟链集群复杂度较高，
> 本毕业设计项目将"区块链引擎"用 Go 语言在 `backend/blockchain/` 目录下
> **完整复现**（工作量证明 + Merkle 树 + 链式结构 + 交易验证）。
> Solidity 智能合约作为**设计稿**保存于本目录，演示了：
> 1. 业务上链的逻辑应该如何用合约表达；
> 2. 部署到 FISCO BCOS 后的预期行为；
> 3. 通过 Solidity 表达 RBAC、状态机、排行榜等关键逻辑的方式。
> 
> 这种"实现 + 合约设计稿"双轨方案既保证系统可运行（Go 链 + MySQL），
> 又完整呈现区块链应用开发的核心能力（Solidity 智能合约）。
