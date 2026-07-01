// SPDX-License-Identifier: MIT
// =============================================================
// 摄影作品投票存证 - 平台管理员合约
pragma solidity ^0.4.25;
import "./Roles.sol";

/**
 * @title Admin
 * @dev 平台管理员角色管理
 *      管理员账户由合约部署者指定，负责：
 *      1. 审核摄影师注册申请
 *      2. 维护黑白名单
 *      3. 处理异常存证请求
 */
contract Admin {
    using Roles for Roles.Role;

    event AdminAdded(address indexed account);
    event AdminRemoved(address indexed account);

    Roles.Role private _admins;

    constructor() public {
        _admins.add(msg.sender);
    }

    modifier onlyAdmin() {
        require(isAdmin(msg.sender), "Admin: caller is not admin");
        _;
    }

    function isAdmin(address account) public view returns (bool) {
        return _admins.has(account);
    }

    function addAdmin(address account) public onlyAdmin {
        _admins.add(account);
        emit AdminAdded(account);
    }

    function renounceAdmin() public {
        _admins.remove(msg.sender);
        emit AdminRemoved(msg.sender);
    }
}
