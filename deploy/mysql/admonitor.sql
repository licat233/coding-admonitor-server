
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for line
-- ----------------------------
DROP TABLE IF EXISTS `line`;
CREATE TABLE `line` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `account` varchar(128) NOT NULL COMMENT '广告账号',
  `cost` decimal(16,6) NOT NULL COMMENT '消耗额',
  `click` int NOT NULL COMMENT '点击数',
  `imp` int NOT NULL COMMENT '展现量',
  `ctr` float NOT NULL COMMENT '点击率',
  `cv` int NOT NULL COMMENT '转化数',
  `cvr` float NOT NULL COMMENT '转化率',
  `cpa` decimal(16,6) NOT NULL COMMENT '转化成本',
  `status` tinyint NOT NULL COMMENT '状态',
  `type` varchar(32) NOT NULL DEFAULT 'newest' COMMENT '数据类型：newest最新数据，node节点数据，day全天数据，logic逻辑数据',
  `currency` varchar(32) NOT NULL DEFAULT 'TWD' COMMENT '币种',
  `record_date` datetime NOT NULL COMMENT '记录日期',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `index_1` (`account`),
  KEY `index_2` (`type`),
  KEY `index_3` (`record_date`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Line直投广告数据';


DROP TABLE IF EXISTS `popin`;
CREATE TABLE `popin` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `account` varchar(128) NOT NULL COMMENT '广告账号',
  `cost` decimal(16,6) NOT NULL COMMENT '消耗额',
  `click` int NOT NULL COMMENT '点击数',
  `imp` int NOT NULL COMMENT '展现量',
  `ctr` float NOT NULL COMMENT '点击率',
  `cv` int NOT NULL COMMENT '转化数',
  `cvr` float NOT NULL COMMENT '转化率',
  `cpa` decimal(16,6) NOT NULL COMMENT '转化成本',
  `status` tinyint NOT NULL COMMENT '状态',
  `type` varchar(32) NOT NULL DEFAULT 'newest' COMMENT '数据类型：newest最新数据，node节点数据，day全天数据，logic逻辑数据',
  `currency` varchar(32) NOT NULL DEFAULT 'TWD' COMMENT '币种',
  `record_date` datetime NOT NULL COMMENT '记录日期',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `index_1` (`account`),
  KEY `index_2` (`type`),
  KEY `index_3` (`record_date`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Line直投广告数据';

SET FOREIGN_KEY_CHECKS = 1;
