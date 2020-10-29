-- ----------------------------
-- 初始化
-- ----------------------------
CREATE TABLE `table1` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `gid` varchar(64) NOT NULL COMMENT '用户id',
  `deleted` tinyint(1) DEFAULT NULL COMMENT '0 未删除，1 已删除',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `tenant_id` varchar(64) DEFAULT NULL,
  `pin` varchar(64) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`pkid`),
  UNIQUE KEY `idx_gid` (`gid`),
  KEY `idx_id` (`id`)
) ENGINE=InnoDB;