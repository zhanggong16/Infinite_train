-- ----------------------------
-- 初始化
-- ----------------------------
CREATE TABLE `instances` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `gid` varchar(36) NOT NULL COMMENT '用户id',
  `uuid` varchar(36) NOT NULL COMMENT '唯一id',
  `name` varchar(36) NOT NULL COMMENT '名称',
  `tenant_id` varchar(64) DEFAULT "",
  `deleted` tinyint(1) DEFAULT 0 COMMENT '0 未删除，1 已删除',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_gid` (`gid`),
  UNIQUE KEY `idx_uuid` (`uuid`)
) ENGINE=InnoDB;