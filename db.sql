/*
 Navicat Premium Data Transfer

 Source Server         : 192.168.1.8
 Source Server Type    : MySQL
 Source Server Version : 80028
 Source Host           : 192.168.1.8:3306
 Source Schema         : fschedule

 Target Server Type    : MySQL
 Target Server Version : 80028
 File Encoding         : 65001

 Date: 13/02/2023 21:31:29
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for task
-- ----------------------------
DROP TABLE IF EXISTS `task`;
CREATE TABLE `task` (
  `Id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '实现Job的特性名称（客户端识别哪个实现类）',
  `ver` int NOT NULL DEFAULT '0' COMMENT '版本',
  `caption` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务组标题',
  `start_at` datetime(6) NOT NULL COMMENT '开始时间',
  `run_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '实际执行时间',
  `run_speed` int NOT NULL COMMENT '运行耗时',
  `client_id` bigint NOT NULL DEFAULT '0' COMMENT '客户端ID',
  `client_ip` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '客户端IP',
  `client_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '客户端名称',
  `progress` int NOT NULL COMMENT '进度0-100',
  `status` tinyint NOT NULL COMMENT '状态',
  `scheduler_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '调度时间',
  `data` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '动态参数',
  `create_at` datetime(6) NOT NULL COMMENT '任务创建时间',
  PRIMARY KEY (`Id`) USING BTREE,
  KEY `group_id_status` (`status`,`create_at`,`Id`) USING BTREE,
  KEY `task_group_id` (`create_at`) USING BTREE,
  KEY `start_at` (`start_at`,`status`) USING BTREE,
  KEY `create_at` (`status`,`create_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=237069 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for task_group
-- ----------------------------
DROP TABLE IF EXISTS `task_group`;
CREATE TABLE `task_group` (
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '实现Job的特性名称（客户端识别哪个实现类）',
  `ver` int NOT NULL DEFAULT '0' COMMENT '版本',
  `caption` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务组标题',
  `start_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '开始时间',
  `next_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '下次执行时间',
  `cron` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '时间定时器表达式',
  `activate_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '活动时间',
  `last_run_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '最后一次完成时间',
  `run_speed_avg` bigint NOT NULL DEFAULT '0' COMMENT '运行平均耗时',
  `run_count` int NOT NULL DEFAULT '0' COMMENT '运行次数',
  `is_enable` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否开启',
  `data` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '动态参数',
  `task` varchar(4096) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务',
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for task_log
-- ----------------------------
DROP TABLE IF EXISTS `task_log`;
CREATE TABLE `task_log` (
  `Id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '实现Job的特性名称（客户端识别哪个实现类）',
  `ver` int NOT NULL DEFAULT '0' COMMENT '版本',
  `caption` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务组标题',
  `task_id` bigint NOT NULL DEFAULT '0' COMMENT '任务记录ID',
  `data` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '动态参数',
  `log_level` tinyint NOT NULL DEFAULT '0' COMMENT '日志级别',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '日志内容',
  `create_at` datetime(6) NOT NULL COMMENT '日志时间',
  PRIMARY KEY (`Id`,`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

SET FOREIGN_KEY_CHECKS = 1;
