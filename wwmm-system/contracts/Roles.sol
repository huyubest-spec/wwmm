// SPDX-License-Identifier: MIT
// =============================================================
// 摄影作品投票存证 - 角色管理基础库
// 该文件提供 RBAC 角色管理功能，可被其他合约通过 import 复用
pragma solidity ^0.4.25;

/**
 * @title Roles
 * @dev 基础角色库：定义一个 Role 类型，并提供 add/remove/has 三个函数
 */
library Roles {
    struct Role {
        mapping(address => bool) bearer;
    }

    function add(Role storage role, address account) internal {
        require(!has(role, account), "Roles: account already has role");
        role.bearer[account] = true;
    }

    function remove(Role storage role, address account) internal {
        require(has(role, account), "Roles: account does not have role");
        role.bearer[account] = false;
    }

    function has(Role storage role, address account) internal view returns (bool) {
        require(account != address(0), "Roles: account is the zero address");
        return role.bearer[account];
    }
}
