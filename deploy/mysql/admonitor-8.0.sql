/*
 Navicat MySQL Data Transfer

 Source Server         : admonitor
 Source Server Type    : MySQL
 Source Server Version : 80024
 Source Host           : 119.91.140.95:3306
 Source Schema         : admonitor

 Target Server Type    : MySQL
 Target Server Version : 80024
 File Encoding         : 65001

 Date: 19/04/2022 18:36:39
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for cookie
-- ----------------------------
DROP TABLE IF EXISTS `cookie`;
CREATE TABLE `cookie` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `channel` varchar(128) NOT NULL COMMENT 'cookie归属渠道，如：line、popin',
  `account` varchar(128) NOT NULL COMMENT 'cookie归属账号，如：licat@gmail.com',
  `name` varchar(128) NOT NULL COMMENT '设置的cookie_name',
  `value` varchar(512) NOT NULL COMMENT '设置的cookie_value',
  `domain` varchar(128) NOT NULL COMMENT 'cookie的domain',
  `path` varchar(128) NOT NULL COMMENT 'cookie的path',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `index_1` (`channel`),
  KEY `index_2` (`account`),
  KEY `index_3` (`name`),
  KEY `index_4` (`domain`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='爬虫cookies表';

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
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Line直投广告数据';

-- ----------------------------
-- Table structure for popin
-- ----------------------------
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
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Line直投广告数据';

-- ----------------------------
-- Table structure for standard
-- ----------------------------
DROP TABLE IF EXISTS `standard`;
CREATE TABLE `standard` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `channel` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '渠道，如line，popin',
  `account` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '账号,如鸡皮肤账号，祛斑账号',
  `imp` int NOT NULL DEFAULT '0' COMMENT '展现',
  `cost` decimal(10,0) NOT NULL DEFAULT '0' COMMENT '消耗',
  `click` int NOT NULL DEFAULT '0' COMMENT '点击',
  `ctr` float NOT NULL DEFAULT '0' COMMENT '点击率',
  `cv` int NOT NULL DEFAULT '0' COMMENT '转化',
  `cvr` float NOT NULL DEFAULT '0' COMMENT '转化率',
  `cpa` decimal(10,0) NOT NULL DEFAULT '0' COMMENT '转化成本',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改日期',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='广告数据标准';

SET FOREIGN_KEY_CHECKS = 1;
