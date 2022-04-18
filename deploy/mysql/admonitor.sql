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

 Date: 18/04/2022 09:53:09
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
) ENGINE=InnoDB AUTO_INCREMENT=86 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='爬虫cookies表';

-- ----------------------------
-- Records of cookie
-- ----------------------------
BEGIN;
INSERT INTO `cookie` VALUES (1, 'popin', 'licat233@gmail.com', '_gid', 'GA1.2.1591965304.1649950445', 'ads.popin.cc', '/', '2022-04-16 16:04:27', '2022-04-16 16:04:27');
INSERT INTO `cookie` VALUES (2, 'popin', 'licat233@gmail.com', 'db_ads', 'Management#assetName%2CAds#headline%2CAds#landingPage%2CAds#brandName%2CCost#spend%2CTrafficAcquisition#cpc%2CTrafficAcquisition#click%2CTrafficAcquisition#ctr%2Cconversion#conversion%2Cconversion#conversionRate', 'ads.popin.cc', '/', '2022-04-16 16:04:27', '2022-04-16 16:04:27');
INSERT INTO `cookie` VALUES (3, 'popin', 'licat233@gmail.com', '_mgruls', 'us', 'ads.popin.cc', '/', '2022-04-16 16:04:27', '2022-04-16 16:04:27');
INSERT INTO `cookie` VALUES (4, 'popin', 'licat233@gmail.com', 'db_campaign', 'Management#campaignName%2CCost#spend%2CCampaigns#dailyBudget%2CTrafficAcquisition#cpc%2CTrafficAcquisition#impression%2CTrafficAcquisition#click%2CTrafficAcquisition#ctr%2Cconversion#conversion%2Cconversion#conversionRate', 'ads.popin.cc', '/', '2022-04-16 16:04:28', '2022-04-16 16:04:28');
INSERT INTO `cookie` VALUES (5, 'popin', 'licat233@gmail.com', 'MEDIAGOOAUTH', 'Ttr71_MNXVMhVCeBzOUs2g6Fu32VMqxNkviyPkIQnq5n6wGoR5iZ4VEk5qYs8ljnPCyWSSZMN6l2So5O_WEP_g', 'ads.popin.cc', '/', '2022-04-16 16:04:28', '2022-04-16 16:04:28');
INSERT INTO `cookie` VALUES (6, 'popin', 'licat233@gmail.com', '__stripe_sid', '656c1903-0af8-4203-886c-66e44e1d6b940ed577', 'ads.popin.cc', '/', '2022-04-16 16:04:28', '2022-04-16 16:04:28');
INSERT INTO `cookie` VALUES (7, 'popin', 'licat233@gmail.com', '__stripe_mid', '41a3a6f0-c2b2-4a20-a1b6-3596fa1235b7fcfbf6', 'ads.popin.cc', '/', '2022-04-16 16:04:28', '2022-04-16 16:04:28');
INSERT INTO `cookie` VALUES (8, 'popin', 'licat233@gmail.com', 'new_db_trend', 'spend', 'ads.popin.cc', '/', '2022-04-16 16:04:28', '2022-04-16 16:04:28');
INSERT INTO `cookie` VALUES (9, 'popin', 'licat233@gmail.com', 'mgvdata', 'JFG96iVCJrcwsrDp', 'ads.popin.cc', '/', '2022-04-16 16:04:28', '2022-04-16 16:04:28');
INSERT INTO `cookie` VALUES (10, 'popin', 'licat233@gmail.com', '_ga_ZX5XQZNMRL', 'GS1.1.1650027777.30.1.1650027778.0', 'ads.popin.cc', '/', '2022-04-16 16:04:28', '2022-04-16 16:04:28');
INSERT INTO `cookie` VALUES (11, 'popin', 'licat233@gmail.com', '_ga', 'GA1.1.1772045387.1649950445', 'ads.popin.cc', '/', '2022-04-16 16:04:28', '2022-04-16 16:04:28');
INSERT INTO `cookie` VALUES (12, 'popin', 'licat233@gmail.com', 'mgrudata', 'id%3D1200%26name%3Dlicat233%2540gmail.com%26type%3Dm%26source%3Dmsnpopin%26t%3D1650020178%26v%3De824f6afca9bde89f1bc465b2eb6cbb7_2f97865a9770d91f36789fd72076c7456204e1b4', 'ads.popin.cc', '/', '2022-04-16 16:04:28', '2022-04-16 16:04:28');
INSERT INTO `cookie` VALUES (13, 'popin', 'licat233@gmail.com', 'mgruli', 'n%3D%26t%3Dmsnpopin', 'ads.popin.cc', '/', '2022-04-16 16:04:28', '2022-04-16 16:04:28');
INSERT INTO `cookie` VALUES (14, 'popin', 'licat233@gmail.com', 'db_account', 'Management#accountName%2CCost#spend%2CTrafficAcquisition#avgCpc%2CTrafficAcquisition#impression%2CTrafficAcquisition#click%2CTrafficAcquisition#ctr', 'ads.popin.cc', '/', '2022-04-16 16:04:28', '2022-04-16 16:04:28');
INSERT INTO `cookie` VALUES (15, 'popin', 'licat233@gmail.com', 'currency', 'TWD', 'ads.popin.cc', '/', '2022-04-16 16:04:28', '2022-04-16 16:04:28');
INSERT INTO `cookie` VALUES (16, 'line', 'licat233@gmail.com', '_trmcuser', '{\"id\":\"\"}', 'line.biz', '/', '2022-04-17 23:06:30', '2022-04-17 23:06:30');
INSERT INTO `cookie` VALUES (17, 'line', 'licat233@gmail.com', '_trmccid', '862056ab30933ed6', 'line.biz', '/', '2022-04-17 23:06:30', '2022-04-17 23:06:30');
INSERT INTO `cookie` VALUES (18, 'line', 'licat233@gmail.com', '_ga', 'GA1.2.570990389.1650125186', 'line.biz', '/', '2022-04-17 23:06:30', '2022-04-17 23:06:30');
INSERT INTO `cookie` VALUES (19, 'line', 'licat233@gmail.com', 'SESSION', 'YzQ3YzJmNTYtNmEzOC00MDM0LWFjNTAtM2E3YjYxNTQ3OTIy', 'line.biz', '/', '2022-04-17 23:06:30', '2022-04-17 23:57:32');
INSERT INTO `cookie` VALUES (20, 'line', 'licat233@gmail.com', '_trmcdisabled2', '-1', 'line.biz', '/', '2022-04-17 23:06:30', '2022-04-17 23:06:30');
INSERT INTO `cookie` VALUES (21, 'line', 'licat233@gmail.com', '__try__', '1650125186350', 'line.biz', '/', '2022-04-17 23:06:30', '2022-04-17 23:06:30');
INSERT INTO `cookie` VALUES (22, 'line', 'licat233@gmail.com', '_gid', 'GA1.2.1883950194.1650125186', 'line.biz', '/', '2022-04-17 23:06:30', '2022-04-17 23:06:30');
INSERT INTO `cookie` VALUES (23, 'line', 'licat233@gmail.com', '_gat_gtag_UA_124047556_1', '1', 'line.biz', '/', '2022-04-17 23:06:31', '2022-04-17 23:06:31');
INSERT INTO `cookie` VALUES (24, 'line', 'licat233@gmail.com', 'XSRF-TOKEN', '9d96d2c1-58cd-4a14-8025-34b15a492733', 'line.biz', '/', '2022-04-17 23:06:31', '2022-04-17 23:57:32');
INSERT INTO `cookie` VALUES (25, 'line', 'licat233@gmail.com', '_trmcsession', '{\"id\":\"f292deba58c17ca0\",\"path\":\"/login\",\"query\":\"?redirectUri+https%3A%2F%2Fadmanager.line.biz%2Flogin%2F&scope+email\",\"params\":{},\"time\":1650125186356}', 'line.biz', '/', '2022-04-17 23:06:31', '2022-04-17 23:06:31');
INSERT INTO `cookie` VALUES (26, 'line', 'licat233@gmail.com', '_trmcpage', '/login', 'line.biz', '/', '2022-04-17 23:06:31', '2022-04-17 23:06:31');
COMMIT;

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
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Line直投广告数据';

-- ----------------------------
-- Records of line
-- ----------------------------
BEGIN;
INSERT INTO `line` VALUES (1, 'A28490221213', 0.090000, 1, 4, 25, 9, 900, 0.010000, 1, 'newest', 'TWD', '2022-04-14 11:13:01', '2022-04-14 11:13:01', '2022-04-14 11:13:01');
INSERT INTO `line` VALUES (2, 'A28490221213', 0.000000, 0, 0, 0, 0, 0, 0.000000, 0, 'node', 'TWD', '2022-04-14 09:30:00', '2022-04-14 11:13:01', '2022-04-14 11:13:01');
INSERT INTO `line` VALUES (3, 'A28490221213', 11541.800000, 2884, 183564, 1.57, 114, 3.95, 101.240000, 1, 'day', 'TWD', '2022-04-13 11:13:02', '2022-04-14 11:13:02', '2022-04-14 11:13:02');
INSERT INTO `line` VALUES (4, 'A28490221213', 0.090000, 1, 4, 25, 9, 900, 0.010000, 0, 'logic', 'TWD', '2022-04-14 09:30:00', '2022-04-14 11:13:02', '2022-04-14 11:13:02');
INSERT INTO `line` VALUES (5, '', 0.000000, 0, 0, 0, 0, 0, 0.000000, 0, 'node', 'TWD', '2022-04-14 09:30:00', '2022-04-14 16:52:54', '2022-04-14 16:52:54');
COMMIT;

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
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Line直投广告数据';

-- ----------------------------
-- Records of popin
-- ----------------------------
BEGIN;
INSERT INTO `popin` VALUES (1, '', 0.000000, 0, 0, 0, 0, 0, 0.000000, 0, 'node', 'TWD', '2022-04-14 09:30:00', '2022-04-14 17:44:46', '2022-04-14 17:44:46');
INSERT INTO `popin` VALUES (2, '', 0.000000, 0, 0, 0, 0, 0, 0.000000, 0, 'node', 'TWD', '2022-04-15 09:30:00', '2022-04-15 16:04:36', '2022-04-15 16:04:36');
INSERT INTO `popin` VALUES (3, 'Wang_skin01', 35220.180000, 7046, 1155957, 0.0061, 0, 0, 0.000000, 1, 'newest', 'TWD', '2022-04-15 22:40:13', '2022-04-15 22:39:56', '2022-04-15 22:40:13');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
