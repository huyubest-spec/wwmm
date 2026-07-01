// SPDX-License-Identifier: MIT
// =============================================================
// 摄影作品投票存证 - 核心合约
pragma solidity ^0.4.25;
import "./Photographer.sol";

/**
 * @title PhotoEvidence
 * @dev 摄影作品存证 + 投票主合约
 *
 * 工作流程：
 * 1. 摄影师注册后，调用 submitPhoto 上传作品（提交文件哈希+元数据）
 * 2. 管理员调用 auditPhoto 审核作品
 * 3. 审核通过后，任何地址均可调用 voteFor 投票，每地址每作品仅 1 票
 * 4. 全部信息永久保存在链上，可通过 getPhoto / getVotes / getRanking 查询
 */
contract PhotoEvidence is Photographer {

    enum PhotoStatus { Pending, Approved, Rejected }

    struct Photo {
        uint256       photoId;
        address       photographer;     // 摄影师地址
        string        title;            // 作品标题
        string        imageHash;        // 图片 SHA-256 哈希（不可篡改）
        string        description;      // 作品描述
        string        category;         // 分类
        string        shootLocation;    // 拍摄地点
        uint256       submitTime;       // 提交时间
        uint256       auditTime;        // 审核时间
        PhotoStatus   status;           // 状态
        uint256       voteCount;        // 累计得票
        string        auditComment;     // 审核意见
        bool          active;           // 是否有效
    }

    struct Vote {
        address voter;
        uint256 photoId;
        uint256 voteTime;
        bool    valid;
    }

    uint256 private _photoSeq;
    mapping(uint256 => Photo) private _photos;
    mapping(uint256 => Vote[]) private _votes;        // photoId => 投票列表
    mapping(address => uint256[]) private _userVotes; // voter => 已投票 photoId 列表
    mapping(bytes32 => uint256) private _hashIndex;   // imageHash => photoId 防止重复

    event PhotoSubmitted(uint256 indexed photoId, address indexed photographer, string imageHash, uint256 submitTime);
    event PhotoAudited(uint256 indexed photoId, PhotoStatus status, string comment);
    event Voted(uint256 indexed photoId, address indexed voter, uint256 voteCount);

    modifier photoExists(uint256 photoId) {
        require(_photos[photoId].photoId != 0, "PhotoEvidence: photo not found");
        _;
    }

    /// 摄影师提交作品存证
    function submitPhoto(
        string memory title,
        string memory imageHash,
        string memory description,
        string memory category,
        string memory shootLocation
    ) public onlyPhotographer returns (uint256) {
        bytes32 h = keccak256(abi.encodePacked(imageHash));
        require(_hashIndex[h] == 0, "PhotoEvidence: photo already submitted");
        require(bytes(title).length > 0, "PhotoEvidence: title required");
        require(bytes(imageHash).length == 64, "PhotoEvidence: image hash must be 64 hex chars");

        _photoSeq += 1;
        uint256 pid = _photoSeq;
        _photos[pid] = Photo({
            photoId:        pid,
            photographer:   msg.sender,
            title:          title,
            imageHash:      imageHash,
            description:    description,
            category:       category,
            shootLocation:  shootLocation,
            submitTime:     now,
            auditTime:      0,
            status:         PhotoStatus.Pending,
            voteCount:      0,
            auditComment:   "",
            active:         true
        });
        _hashIndex[h] = pid;
        _incrementPhotoCount(msg.sender);

        emit PhotoSubmitted(pid, msg.sender, imageHash, now);
        return pid;
    }

    /// 管理员审核作品
    function auditPhoto(uint256 photoId, bool approve, string memory comment) public onlyAdmin photoExists(photoId) {
        Photo storage p = _photos[photoId];
        require(p.status == PhotoStatus.Pending, "PhotoEvidence: photo already audited");
        p.status = approve ? PhotoStatus.Approved : PhotoStatus.Rejected;
        p.auditComment = comment;
        p.auditTime = now;
        emit PhotoAudited(photoId, p.status, comment);
    }

    /// 投票
    function voteFor(uint256 photoId) public photoExists(photoId) {
        Photo storage p = _photos[photoId];
        require(p.status == PhotoStatus.Approved, "PhotoEvidence: photo not approved");
        require(p.photographer != msg.sender, "PhotoEvidence: cannot vote for yourself");
        require(_userVotes[msg.sender].length == 0 || !_hasVoted(msg.sender, photoId),
                "PhotoEvidence: already voted");

        _votes[photoId].push(Vote({
            voter:    msg.sender,
            photoId:  photoId,
            voteTime: now,
            valid:    true
        }));
        _userVotes[msg.sender].push(photoId);
        p.voteCount += 1;
        emit Voted(photoId, msg.sender, p.voteCount);
    }

    function _hasVoted(address voter, uint256 photoId) internal view returns (bool) {
        uint256[] memory list = _userVotes[voter];
        for (uint256 i = 0; i < list.length; i++) {
            if (list[i] == photoId) return true;
        }
        return false;
    }

    /// 查询作品
    function getPhoto(uint256 photoId) public view photoExists(photoId) returns (
        uint256 pid,
        address photographer,
        string memory title,
        string memory imageHash,
        string memory description,
        uint256 submitTime,
        uint256 auditTime,
        PhotoStatus status,
        uint256 voteCount,
        string memory auditComment
    ) {
        Photo storage p = _photos[photoId];
        return (
            p.photoId, p.photographer, p.title, p.imageHash, p.description,
            p.submitTime, p.auditTime, p.status, p.voteCount, p.auditComment
        );
    }

    /// 通过图片哈希反查作品
    function getPhotoByHash(string memory imageHash) public view returns (uint256) {
        bytes32 h = keccak256(abi.encodePacked(imageHash));
        return _hashIndex[h];
    }

    /// 查询投票记录
    function getVotes(uint256 photoId) public view photoExists(photoId) returns (uint256) {
        return _votes[photoId].length;
    }

    /// 查询排行榜 (取前 N 个通过审核的作品，按得票降序)
    function getRanking(uint256 topN) public view returns (uint256[] memory) {
        uint256[] memory tmp = new uint256[](_photoSeq);
        uint256 n = 0;
        for (uint256 i = 1; i <= _photoSeq; i++) {
            if (_photos[i].status == PhotoStatus.Approved && _photos[i].active) {
                tmp[n++] = i;
            }
        }
        // 简化版：直接冒泡降序排列前 n 个
        for (uint256 i = 0; i < n; i++) {
            for (uint256 j = i + 1; j < n; j++) {
                if (_photos[tmp[j]].voteCount > _photos[tmp[i]].voteCount) {
                    uint256 t = tmp[i];
                    tmp[i] = tmp[j];
                    tmp[j] = t;
                }
            }
        }
        uint256 len = n < topN ? n : topN;
        uint256[] memory result = new uint256[](len);
        for (uint256 i = 0; i < len; i++) {
            result[i] = tmp[i];
        }
        return result;
    }

    function photoCount() public view returns (uint256) {
        return _photoSeq;
    }
}
