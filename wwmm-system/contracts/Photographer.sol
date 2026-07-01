// SPDX-License-Identifier: MIT
// =============================================================
// 摄影作品投票存证 - 摄影师合约
pragma solidity ^0.4.25;
import "./Roles.sol";
import "./Admin.sol";

/**
 * @title Photographer
 * @dev 摄影师角色管理
 *      摄影师可以：
 *      1. 上传作品（提交存证请求）
 *      2. 撤回未审核的作品
 *      3. 查看自己的作品列表
 */
contract Photographer is Admin {
    using Roles for Roles.Role;

    struct PhotographerInfo {
        string  realName;
        string  phone;
        uint256 registerTime;
        bool    active;
        uint256 photoCount;
    }

    event PhotographerRegistered(address indexed account, string realName);
    event PhotographerRemoved(address indexed account);

    Roles.Role private _photographers;
    mapping(address => PhotographerInfo) private _info;

    modifier onlyPhotographer() {
        require(isPhotographer(msg.sender), "Photographer: caller is not photographer");
        _;
    }

    function isPhotographer(address account) public view returns (bool) {
        return _photographers.has(account);
    }

    function register(string memory realName, string memory phone) public {
        require(!isPhotographer(msg.sender), "Photographer: already registered");
        _photographers.add(msg.sender);
        _info[msg.sender] = PhotographerInfo({
            realName:     realName,
            phone:        phone,
            registerTime: now,
            active:       true,
            photoCount:   0
        });
        emit PhotographerRegistered(msg.sender, realName);
    }

    function removePhotographer(address account) public onlyAdmin {
        _photographers.remove(account);
        _info[account].active = false;
        emit PhotographerRemoved(account);
    }

    function getInfo(address account) public view returns (
        string memory realName,
        string memory phone,
        uint256 registerTime,
        bool    active,
        uint256 photoCount
    ) {
        PhotographerInfo storage info = _info[account];
        return (info.realName, info.phone, info.registerTime, info.active, info.photoCount);
    }

    function _incrementPhotoCount(address account) internal {
        _info[account].photoCount += 1;
    }
}
