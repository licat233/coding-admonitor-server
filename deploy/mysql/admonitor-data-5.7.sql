/*
 Navicat MySQL Data Transfer

 Source Server         : old-bt-admonitor
 Source Server Type    : MySQL
 Source Server Version : 50734
 Source Host           : 150.109.61.90:3306
 Source Schema         : admonitor

 Target Server Type    : MySQL
 Target Server Version : 50734
 File Encoding         : 65001

 Date: 02/05/2022 12:08:24
*/

SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for account
-- ----------------------------
DROP TABLE IF EXISTS `account`;
CREATE TABLE `account` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '表主键',
  `aid` varchar(128) NOT NULL COMMENT '账号id',
  `name` varchar(128) NOT NULL COMMENT '账号名',
  `nickname` varchar(128) NOT NULL COMMENT '账号昵称',
  `channel_id` int(11) NOT NULL COMMENT '归属渠道channel表主键id',
  `status` tinyint(4) NOT NULL COMMENT '账号状态，0关闭,1开启',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `aid` (`aid`,`channel_id`,`status`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8 COMMENT='广告子账号表，并非登录账户';

-- ----------------------------
-- Records of account
-- ----------------------------
BEGIN;
INSERT INTO `account` VALUES (1, '1601', 'Wang_skin01', '鸡皮肤1户', 1, 1, '2022-04-29 15:06:37', '2022-04-29 15:06:37');
INSERT INTO `account` VALUES (2, '1445', 'Self_White3', '祛斑3户', 1, 0, '2022-04-29 15:06:37', '2022-04-29 22:25:28');
INSERT INTO `account` VALUES (3, '1442', 'Self_White2', '祛斑2户', 1, 0, '2022-04-29 15:06:37', '2022-04-29 15:06:37');
INSERT INTO `account` VALUES (4, '1540', 'Wang_skin02', '鸡皮肤2户', 1, 0, '2022-04-29 15:06:37', '2022-04-29 15:06:37');
INSERT INTO `account` VALUES (5, '1618', 'Wang_OLD', '抗老1户', 1, 0, '2022-04-29 15:06:37', '2022-04-29 15:06:37');
INSERT INTO `account` VALUES (6, '16826', 'Self_wrinkle', '抗老2户', 1, 0, '2022-04-29 15:06:37', '2022-04-29 15:06:37');
INSERT INTO `account` VALUES (7, 'A79362240622', '思原生物工程有限公司 - 拾光肌因', '鸡皮肤账户', 2, 1, '2022-04-29 15:20:15', '2022-04-29 16:08:00');
INSERT INTO `account` VALUES (8, 'A28490221213', '	 艾凱創意行銷股份有限公司 - 植題2', '祛斑账户', 2, 1, '2022-04-29 15:20:15', '2022-04-29 16:07:52');
INSERT INTO `account` VALUES (9, 'A96457201158', '艾凱創意行銷股份有限公司 - 鮮活煥顏精華原液4', '被冻结账户', 2, 0, '2022-04-29 15:20:15', '2022-04-29 16:08:12');
INSERT INTO `account` VALUES (10, 'A97230191467', '艾凱創意行銷股份有限公司 - 鮮活煥顏精華原液2', '被冻结账户', 2, 0, '2022-04-29 15:20:15', '2022-04-29 16:08:19');
COMMIT;

-- ----------------------------
-- Table structure for channel
-- ----------------------------
DROP TABLE IF EXISTS `channel`;
CREATE TABLE `channel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL COMMENT '渠道名称',
  `email` varchar(128) NOT NULL COMMENT '渠道登录邮箱',
  `passport` varchar(128) NOT NULL COMMENT '渠道登录账号',
  `password` varchar(128) NOT NULL COMMENT '渠道登录密码',
  `domain` varchar(128) NOT NULL COMMENT '渠道cookie根域，例如: .line.biz',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '渠道状态,0关闭，1开启',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='渠道表';

-- ----------------------------
-- Records of channel
-- ----------------------------
BEGIN;
INSERT INTO `channel` VALUES (1, 'popin', 'licat233@gmail.com', '', '205D7E435CD6230A13279F251ED4B3C4', 'ads.popin.cc', 1, '2022-04-29 11:44:53', '2022-04-29 14:48:48');
INSERT INTO `channel` VALUES (2, 'line', 'licat233@gmail.com', '', 'wang412line...', '.line.biz', 0, '2022-04-29 11:45:36', '2022-04-29 22:23:42');
COMMIT;

-- ----------------------------
-- Table structure for cookie
-- ----------------------------
DROP TABLE IF EXISTS `cookie`;
CREATE TABLE `cookie` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `channel_id` int(11) NOT NULL COMMENT 'cookie归属渠道channel表主键id',
  `name` varchar(128) NOT NULL COMMENT '设置的cookie_name',
  `value` varchar(512) NOT NULL COMMENT '设置的cookie_value',
  `domain` varchar(128) NOT NULL COMMENT 'cookie的domain',
  `path` varchar(128) NOT NULL COMMENT 'cookie的path',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `index_1` (`channel_id`),
  KEY `index_3` (`name`),
  KEY `index_4` (`domain`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='爬虫cookies表';

-- ----------------------------
-- Table structure for standard
-- ----------------------------
DROP TABLE IF EXISTS `standard`;
CREATE TABLE `standard` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `account_id` int(11) NOT NULL COMMENT '归属渠道channel表主键id',
  `imp` int(11) NOT NULL DEFAULT '0' COMMENT '展现',
  `cost` decimal(16,6) NOT NULL DEFAULT '0.000000' COMMENT '消耗',
  `click` int(11) NOT NULL DEFAULT '0' COMMENT '点击',
  `ctr` float NOT NULL DEFAULT '0' COMMENT '点击率',
  `cv` int(11) NOT NULL DEFAULT '0' COMMENT '转化',
  `cvr` float NOT NULL DEFAULT '0' COMMENT '转化率',
  `cpa` decimal(16,6) NOT NULL DEFAULT '0.000000' COMMENT '转化成本',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '监控状态,0不监控，1监控',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
  PRIMARY KEY (`id`),
  KEY `index_2` (`account_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='广告数据规范';

-- ----------------------------
-- Records of standard
-- ----------------------------
BEGIN;
INSERT INTO `standard` VALUES (1, 2, 0, 0.000000, 0, 0, 0, 0, 0.000000, 0, '2022-04-29 16:10:03', '2022-04-29 16:10:03');
INSERT INTO `standard` VALUES (2, 1, 0, 0.000000, 0, 0, 0, 0, 0.000000, 0, '2022-04-29 16:10:03', '2022-04-29 16:10:03');
COMMIT;

-- ----------------------------
-- Table structure for statistics
-- ----------------------------
DROP TABLE IF EXISTS `statistics`;
CREATE TABLE `statistics` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `identifier` varchar(255) NOT NULL COMMENT '唯一标识符，由channel_id、account_id、record_date、type用-符号连接而成，目的是为了防止重复插入',
  `channel_id` int(11) NOT NULL COMMENT '渠道id，指的是channel表的主键',
  `account_id` int(11) NOT NULL COMMENT '账号id，指的是account表中的主键',
  `cost` decimal(16,6) NOT NULL COMMENT '消耗额',
  `click` int(11) NOT NULL COMMENT '点击数',
  `imp` int(11) NOT NULL COMMENT '展现量',
  `ctr` float NOT NULL COMMENT '点击率',
  `cv` int(11) NOT NULL COMMENT '转化数',
  `cvr` float NOT NULL COMMENT '转化率',
  `cpa` decimal(16,6) NOT NULL COMMENT '转化成本',
  `status` tinyint(4) NOT NULL COMMENT '状态',
  `type` varchar(32) NOT NULL DEFAULT 'newest' COMMENT '数据类型：newest最新数据，node节点数据，day全天数据，logic逻辑数据',
  `currency` varchar(32) NOT NULL DEFAULT 'TWD' COMMENT '币种',
  `record_date` datetime NOT NULL COMMENT '记录日期，09:30-09:30为一天',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `identifier` (`identifier`),
  KEY `index_1` (`channel_id`),
  KEY `index_2` (`account_id`),
  KEY `index_3` (`type`),
  KEY `index_4` (`record_date`),
  KEY `identifier_2` (`identifier`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='广告投放监控数据';

SET FOREIGN_KEY_CHECKS = 1;
