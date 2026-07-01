-- ============================================================
--  基于区块链的摄影作品投票存证系统
--  数据库初始化脚本
--  适用于 MySQL 8.0+
-- ============================================================

DROP DATABASE IF EXISTS wwmm_db;
CREATE DATABASE wwmm_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE wwmm_db;

-- ============================================================
--  1. 用户表 user
--  角色说明：0-普通用户 1-摄影师 2-管理员
-- ============================================================
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `user_id`        INT          NOT NULL AUTO_INCREMENT         COMMENT '用户ID',
  `username`       VARCHAR(64)  NOT NULL                         COMMENT '登录账号',
  `password_hash`  VARCHAR(255) NOT NULL                         COMMENT '密码哈希(SHA-256+盐)',
  `salt`           VARCHAR(32)  NOT NULL                         COMMENT '盐值',
  `phone`          VARCHAR(20)  DEFAULT NULL                     COMMENT '手机号',
  `email`          VARCHAR(128) DEFAULT NULL                     COMMENT '邮箱',
  `real_name`      VARCHAR(64)  DEFAULT NULL                     COMMENT '真实姓名',
  `sex`            TINYINT      DEFAULT 0                        COMMENT '性别 0-未知 1-男 2-女',
  `avatar`         VARCHAR(255) DEFAULT NULL                     COMMENT '头像URL',
  `bio`            VARCHAR(255) DEFAULT NULL                     COMMENT '个人简介',
  `role`           TINYINT      NOT NULL DEFAULT 0               COMMENT '角色 0-普通用户 1-摄影师 2-管理员',
  `status`         TINYINT      NOT NULL DEFAULT 1               COMMENT '账号状态 0-禁用 1-启用',
  `created_at`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_role` (`role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- ============================================================
--  2. 摄影作品表 photo
--  状态说明：0-待审核 1-已通过 2-已拒绝
-- ============================================================
DROP TABLE IF EXISTS `photo`;
CREATE TABLE `photo` (
  `photo_id`        INT          NOT NULL AUTO_INCREMENT         COMMENT '作品ID',
  `title`           VARCHAR(128) NOT NULL                         COMMENT '作品标题',
  `description`     TEXT         DEFAULT NULL                     COMMENT '作品描述',
  `image_url`       VARCHAR(255) NOT NULL                         COMMENT '图片访问URL',
  `image_hash`      VARCHAR(64)  NOT NULL                         COMMENT '图片SHA-256哈希(上链存证)',
  `file_size`       INT          DEFAULT 0                        COMMENT '文件大小(字节)',
  `photographer_id` INT          NOT NULL                         COMMENT '摄影师ID',
  `category`        VARCHAR(32)  DEFAULT NULL                     COMMENT '分类(风光/人文/纪实等)',
  `shoot_location`  VARCHAR(128) DEFAULT NULL                     COMMENT '拍摄地点',
  `shoot_time`      DATE         DEFAULT NULL                     COMMENT '拍摄时间',
  `camera_info`     VARCHAR(128) DEFAULT NULL                     COMMENT '拍摄器材信息',
  `status`          TINYINT      NOT NULL DEFAULT 0               COMMENT '状态 0-待审核 1-已通过 2-已拒绝',
  `audit_comment`   VARCHAR(255) DEFAULT NULL                     COMMENT '审核意见',
  `vote_count`      INT          NOT NULL DEFAULT 0               COMMENT '得票数(缓存)',
  `view_count`      INT          NOT NULL DEFAULT 0               COMMENT '浏览数',
  `is_on_chain`     TINYINT      NOT NULL DEFAULT 0               COMMENT '是否已上链 0-否 1-是',
  `chain_tx_hash`   VARCHAR(64)  DEFAULT NULL                     COMMENT '上链交易哈希',
  `created_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`photo_id`),
  UNIQUE KEY `uk_image_hash` (`image_hash`),
  KEY `idx_photographer` (`photographer_id`),
  KEY `idx_status` (`status`),
  KEY `idx_vote_count` (`vote_count` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='摄影作品表';

-- ============================================================
--  3. 投票记录表 vote
--  每个用户对每个作品只能投 1 票
-- ============================================================
DROP TABLE IF EXISTS `vote`;
CREATE TABLE `vote` (
  `vote_id`     INT      NOT NULL AUTO_INCREMENT                    COMMENT '投票ID',
  `user_id`     INT      NOT NULL                                    COMMENT '投票用户ID',
  `photo_id`    INT      NOT NULL                                    COMMENT '作品ID',
  `tx_hash`     VARCHAR(64) DEFAULT NULL                             COMMENT '投票上链交易哈希',
  `created_at`  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP           COMMENT '投票时间',
  PRIMARY KEY (`vote_id`),
  UNIQUE KEY `uk_user_photo` (`user_id`, `photo_id`),
  KEY `idx_photo` (`photo_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='投票记录表';

-- ============================================================
--  4. 区块表 block
--  区块链核心表
-- ============================================================
DROP TABLE IF EXISTS `block`;
CREATE TABLE `block` (
  `block_id`     INT          NOT NULL AUTO_INCREMENT         COMMENT '区块ID',
  `index_num`    INT          NOT NULL                         COMMENT '区块高度(从1开始)',
  `prev_hash`    VARCHAR(64)  NOT NULL                         COMMENT '前一区块哈希',
  `merkle_root`  VARCHAR(64)  NOT NULL                         COMMENT 'Merkle根',
  `timestamp`    BIGINT       NOT NULL                         COMMENT '时间戳(秒)',
  `nonce`        BIGINT       NOT NULL DEFAULT 0               COMMENT '工作量证明随机数',
  `difficulty`   INT          NOT NULL DEFAULT 4               COMMENT '挖矿难度',
  `hash`         VARCHAR(64)  NOT NULL                         COMMENT '本区块哈希',
  `tx_count`     INT          NOT NULL DEFAULT 0               COMMENT '交易数',
  `miner`        VARCHAR(128) DEFAULT NULL                     COMMENT '打包者',
  `created_at`   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '入库时间',
  PRIMARY KEY (`block_id`),
  UNIQUE KEY `uk_index` (`index_num`),
  UNIQUE KEY `uk_hash` (`hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='区块表';

-- ============================================================
--  5. 交易表 transaction
--  类型：1-作品存证  2-投票存证
-- ============================================================
DROP TABLE IF EXISTS `tx`;
CREATE TABLE `tx` (
  `tx_id`      INT          NOT NULL AUTO_INCREMENT         COMMENT '交易ID',
  `tx_hash`    VARCHAR(64)  NOT NULL                         COMMENT '交易哈希',
  `block_id`   INT          DEFAULT NULL                     COMMENT '所在区块ID(未打包前为空)',
  `tx_type`    TINYINT      NOT NULL                         COMMENT '交易类型 1-作品存证 2-投票存证',
  `sender`     VARCHAR(128) NOT NULL                         COMMENT '发送方(用户账号)',
  `payload`    TEXT         NOT NULL                         COMMENT '交易载荷(JSON)',
  `status`     TINYINT      NOT NULL DEFAULT 0               COMMENT '状态 0-待打包 1-已打包 2-失败',
  `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`tx_id`),
  UNIQUE KEY `uk_tx_hash` (`tx_hash`),
  KEY `idx_block` (`block_id`),
  KEY `idx_status` (`status`),
  KEY `idx_type` (`tx_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交易表';

-- ============================================================
--  6. 作品审核日志 photo_audit_log
-- ============================================================
DROP TABLE IF EXISTS `photo_audit_log`;
CREATE TABLE `photo_audit_log` (
  `log_id`     INT          NOT NULL AUTO_INCREMENT         COMMENT '日志ID',
  `photo_id`   INT          NOT NULL                         COMMENT '作品ID',
  `admin_id`   INT          NOT NULL                         COMMENT '审核管理员ID',
  `action`     VARCHAR(32)  NOT NULL                         COMMENT '动作(APPROVE/REJECT)',
  `comment`    VARCHAR(255) DEFAULT NULL                     COMMENT '审核意见',
  `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '审核时间',
  PRIMARY KEY (`log_id`),
  KEY `idx_photo` (`photo_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='作品审核日志';

-- ============================================================
--  7. 区块链状态表 chain_state
--  保存当前最新区块索引
-- ============================================================
DROP TABLE IF EXISTS `chain_state`;
CREATE TABLE `chain_state` (
  `id`               INT NOT NULL DEFAULT 1                    COMMENT '固定ID=1',
  `latest_index`     INT NOT NULL DEFAULT 0                    COMMENT '最新区块高度',
  `latest_hash`      VARCHAR(64) DEFAULT NULL                  COMMENT '最新区块哈希',
  `total_tx`         INT NOT NULL DEFAULT 0                    COMMENT '总交易数',
  `updated_at`       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='区块链状态';
INSERT INTO `chain_state` (`id`, `latest_index`, `latest_hash`, `total_tx`) VALUES (1, 0, NULL, 0);

-- ============================================================
--  8. 会话表 session  (用于登录状态保持)
-- ============================================================
DROP TABLE IF EXISTS `session`;
CREATE TABLE `session` (
  `session_id`  VARCHAR(64)  NOT NULL                         COMMENT '会话ID',
  `user_id`     INT          NOT NULL                         COMMENT '用户ID',
  `token`       VARCHAR(255) NOT NULL                         COMMENT 'Token',
  `expire_at`   DATETIME     NOT NULL                         COMMENT '过期时间',
  `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`session_id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会话表';

-- ============================================================
--  初始数据
--  默认密码 admin/123456, photographer/123456, voter/123456
--  实际入库时由后端用 SHA-256(password+salt) 重新生成
-- ============================================================
INSERT INTO `user` (`username`, `password_hash`, `salt`, `phone`, `email`, `real_name`, `role`, `bio`) VALUES
('admin',         'PRECOMPUTED', 'PLACEHOLDER', '13800000001', 'admin@wwmm.com',         '系统管理员', 2, '平台超级管理员'),
('photographer',  'PRECOMPUTED', 'PLACEHOLDER', '13800000002', 'photo@wwmm.com',         '张摄影',     1, '职业风光摄影师'),
('voter',         'PRECOMPUTED', 'PLACEHOLDER', '13800000003', 'voter@wwmm.com',         '王观众',     0, '摄影爱好者');

-- ============================================================
-- 视图 v_photo_full (作品完整信息)
-- ============================================================
CREATE OR REPLACE VIEW v_photo_full AS
SELECT
  p.photo_id, p.title, p.description, p.image_url, p.image_hash, p.file_size,
  p.category, p.shoot_location, p.shoot_time, p.camera_info,
  p.status, p.vote_count, p.view_count, p.is_on_chain, p.chain_tx_hash,
  p.created_at, p.updated_at,
  p.photographer_id, u.username AS photographer_name, u.avatar AS photographer_avatar
FROM photo p LEFT JOIN `user` u ON p.photographer_id = u.user_id;

-- ============================================================
-- 视图 v_block_with_tx (区块及交易概览)
-- ============================================================
CREATE OR REPLACE VIEW v_block_with_tx AS
SELECT
  b.block_id, b.index_num, b.prev_hash, b.merkle_root, b.timestamp, b.nonce, b.difficulty,
  b.hash, b.tx_count, b.miner, b.created_at,
  (SELECT COUNT(*) FROM tx t WHERE t.block_id = b.block_id AND t.tx_type = 1) AS certify_count,
  (SELECT COUNT(*) FROM tx t WHERE t.block_id = b.block_id AND t.tx_type = 2) AS vote_tx_count
FROM block b;
