---
title: Database Models
keywords: component
tags: [component]
sidebar: home_sidebar
permalink: database-models.html
summary: V1 Specification
---

## Database Models

```
/*
Navicat MySQL Data Transfer

Source Server         : 104.196.228.46test
Source Server Version : 50716
Source Host           : 104.196.228.46:3306
Source Database       : newdb

Target Server Type    : MYSQL
Target Server Version : 50716
File Encoding         : 65001

Date: 2016-12-28 12:09:38
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for action
-- ----------------------------
DROP TABLE IF EXISTS `action`;
CREATE TABLE `action` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `namespace` varchar(255) NOT NULL,
  `repository` varchar(255) NOT NULL,
  `workflow` bigint(20) NOT NULL DEFAULT '0',
  `stage` bigint(20) NOT NULL DEFAULT '0',
  `component` bigint(20) NOT NULL DEFAULT '0',
  `service` bigint(20) NOT NULL DEFAULT '0',
  `action` varchar(255) NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` text,
  `event` bigint(20) DEFAULT '0',
  `manifest` longtext,
  `environment` text,
  `kubernetes` text,
  `swarm` text,
  `input` text,
  `output` text,
  `image_name` varchar(255) DEFAULT NULL,
  `image_tag` varchar(255) DEFAULT NULL,
  `timeout` varchar(255) DEFAULT NULL,
  `requires` longtext,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_action_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=1839 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for action_log
-- ----------------------------
DROP TABLE IF EXISTS `action_log`;
CREATE TABLE `action_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `namespace` varchar(255) NOT NULL,
  `repository` varchar(255) NOT NULL,
  `workflow` bigint(20) NOT NULL DEFAULT '0',
  `from_workflow` bigint(20) NOT NULL DEFAULT '0',
  `sequence` bigint(20) NOT NULL DEFAULT '0',
  `stage` bigint(20) NOT NULL DEFAULT '0',
  `from_stage` bigint(20) NOT NULL DEFAULT '0',
  `from_action` bigint(20) NOT NULL DEFAULT '0',
  `run_state` bigint(20) DEFAULT NULL,
  `component` bigint(20) NOT NULL DEFAULT '0',
  `service` bigint(20) NOT NULL DEFAULT '0',
  `action` varchar(255) NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` text,
  `event` bigint(20) DEFAULT '0',
  `manifest` longtext,
  `environment` text,
  `kubernetes` text,
  `swarm` text,
  `input` text,
  `output` text,
  `image_name` varchar(255) DEFAULT NULL,
  `image_tag` varchar(255) DEFAULT NULL,
  `timeout` varchar(255) DEFAULT NULL,
  `requires` longtext,
  `auth_list` longtext,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `fail_reason` varchar(255) DEFAULT NULL,
  `container_id` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_action_log_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2357 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for component
-- ----------------------------
DROP TABLE IF EXISTS `component`;
CREATE TABLE `component` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `namespace` varchar(255) NOT NULL,
  `version` text,
  `version_code` bigint(20) DEFAULT NULL,
  `component` varchar(255) NOT NULL,
  `type` bigint(20) NOT NULL DEFAULT '0',
  `title` varchar(255) DEFAULT NULL,
  `gravatar` text,
  `description` text,
  `endpoint` text,
  `source` text NOT NULL,
  `environment` text,
  `tag` varchar(255) DEFAULT NULL,
  `volume_location` text,
  `volume_data` text,
  `makefile` text,
  `kubernetes` text,
  `swarm` text,
  `input` text,
  `output` text,
  `timeout` bigint(20) DEFAULT NULL,
  `manifest` longtext,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_component_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=66 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for event
-- ----------------------------
DROP TABLE IF EXISTS `event`;
CREATE TABLE `event` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `definition` bigint(20) NOT NULL DEFAULT '0',
  `title` varchar(255) NOT NULL,
  `header` text NOT NULL,
  `payload` longtext NOT NULL,
  `authorization` text,
  `type` bigint(20) NOT NULL DEFAULT '0',
  `source` bigint(20) NOT NULL DEFAULT '0',
  `character` bigint(20) NOT NULL DEFAULT '0',
  `namespace` varchar(255) DEFAULT NULL,
  `repository` varchar(255) DEFAULT NULL,
  `workflow` bigint(20) NOT NULL DEFAULT '0',
  `stage` bigint(20) NOT NULL DEFAULT '0',
  `action` bigint(20) NOT NULL DEFAULT '0',
  `sequence` bigint(20) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_event_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=21561 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for event_definition
-- ----------------------------
DROP TABLE IF EXISTS `event_definition`;
CREATE TABLE `event_definition` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `event` varchar(255) NOT NULL,
  `title` varchar(255) DEFAULT NULL,
  `namespace` varchar(255) DEFAULT NULL,
  `repository` varchar(255) DEFAULT NULL,
  `workflow` bigint(20) DEFAULT NULL,
  `stage` bigint(20) DEFAULT NULL,
  `action` bigint(20) NOT NULL DEFAULT '0',
  `character` bigint(20) NOT NULL DEFAULT '0',
  `type` bigint(20) NOT NULL DEFAULT '0',
  `source` bigint(20) NOT NULL DEFAULT '0',
  `definition` text,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_event_definition_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=14137 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for event_json
-- ----------------------------
DROP TABLE IF EXISTS `event_json`;
CREATE TABLE `event_json` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `site` varchar(255) DEFAULT NULL,
  `type` varchar(255) DEFAULT NULL,
  `output` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for outcome
-- ----------------------------
DROP TABLE IF EXISTS `outcome`;
CREATE TABLE `outcome` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `workflow` bigint(20) NOT NULL DEFAULT '0',
  `real_workflow` bigint(20) NOT NULL DEFAULT '0',
  `stage` bigint(20) NOT NULL DEFAULT '0',
  `real_stage` bigint(20) NOT NULL DEFAULT '0',
  `action` bigint(20) NOT NULL DEFAULT '0',
  `real_action` bigint(20) NOT NULL DEFAULT '0',
  `event` bigint(20) DEFAULT '0',
  `sequence` bigint(20) NOT NULL DEFAULT '0',
  `status` tinyint(1) DEFAULT NULL,
  `result` longtext,
  `output` longtext,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_outcome_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3357 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for serivce_definition
-- ----------------------------
DROP TABLE IF EXISTS `serivce_definition`;
CREATE TABLE `serivce_definition` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `service` varchar(255) NOT NULL,
  `title` varchar(255) DEFAULT NULL,
  `gravatar` text,
  `endpoints` text,
  `status` text,
  `environment` text,
  `authorization` text,
  `configuration` text,
  `description` text,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `service` (`service`),
  KEY `idx_serivce_definition_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for service
-- ----------------------------
DROP TABLE IF EXISTS `service`;
CREATE TABLE `service` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `namespace` varchar(255) NOT NULL,
  `service` varchar(255) NOT NULL,
  `title` varchar(255) DEFAULT NULL,
  `gravatar` text,
  `endpoints` text,
  `environment` text,
  `authorization` text,
  `configuration` text,
  `description` text,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `namespace_service` (`namespace`,`service`),
  KEY `idx_service_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for stage
-- ----------------------------
DROP TABLE IF EXISTS `stage`;
CREATE TABLE `stage` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `namespace` varchar(255) NOT NULL,
  `repository` varchar(255) NOT NULL,
  `workflow` bigint(20) NOT NULL DEFAULT '0',
  `type` bigint(20) NOT NULL DEFAULT '0',
  `pre_stage` bigint(20) NOT NULL DEFAULT '0',
  `stage` varchar(255) NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` text,
  `event` bigint(20) DEFAULT '0',
  `manifest` longtext,
  `env` longtext,
  `timeout` varchar(255) DEFAULT NULL,
  `requires` longtext,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_stage_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=1836 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for stage_log
-- ----------------------------
DROP TABLE IF EXISTS `stage_log`;
CREATE TABLE `stage_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `namespace` varchar(255) NOT NULL,
  `repository` varchar(255) NOT NULL,
  `workflow` bigint(20) NOT NULL DEFAULT '0',
  `from_workflow` bigint(20) NOT NULL DEFAULT '0',
  `sequence` bigint(20) NOT NULL DEFAULT '0',
  `from_stage` bigint(20) NOT NULL DEFAULT '0',
  `type` bigint(20) NOT NULL DEFAULT '0',
  `pre_stage` bigint(20) NOT NULL DEFAULT '0',
  `stage` varchar(255) NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` text,
  `run_state` bigint(20) DEFAULT NULL,
  `event` bigint(20) DEFAULT '0',
  `manifest` longtext,
  `env` longtext,
  `timeout` varchar(255) DEFAULT NULL,
  `requires` longtext,
  `auth_list` longtext,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `fail_reason` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_stage_log_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2306 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for timer
-- ----------------------------
DROP TABLE IF EXISTS `timer`;
CREATE TABLE `timer` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `namespace` varchar(255) NOT NULL,
  `repository` varchar(255) NOT NULL,
  `workflow` bigint(20) DEFAULT NULL,
  `available` tinyint(1) DEFAULT NULL,
  `cron` varchar(255) DEFAULT NULL,
  `event_type` varchar(255) DEFAULT NULL,
  `event_name` varchar(255) DEFAULT NULL,
  `start_json` longtext,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_timer_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user_setting
-- ----------------------------
DROP TABLE IF EXISTS `user_setting`;
CREATE TABLE `user_setting` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `namespace` varchar(255) NOT NULL,
  `repository` varchar(255) NOT NULL,
  `setting` text,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_setting_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for workflow
-- ----------------------------
DROP TABLE IF EXISTS `workflow`;
CREATE TABLE `workflow` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `namespace` varchar(255) NOT NULL,
  `repository` varchar(255) NOT NULL,
  `workflow` varchar(255) NOT NULL,
  `event` bigint(20) DEFAULT '0',
  `version` varchar(255) DEFAULT NULL,
  `version_code` varchar(255) DEFAULT NULL,
  `state` bigint(20) DEFAULT NULL,
  `manifest` longtext,
  `description` text,
  `source_info` longtext,
  `env` longtext,
  `requires` longtext,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `is_limit_instance` tinyint(1) DEFAULT NULL,
  `limit_instance` bigint(20) DEFAULT NULL,
  `current_instance` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_workflow_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=56 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for workflow_log
-- ----------------------------
DROP TABLE IF EXISTS `workflow_log`;
CREATE TABLE `workflow_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `namespace` varchar(255) NOT NULL,
  `repository` varchar(255) NOT NULL,
  `workflow` varchar(255) NOT NULL,
  `from_workflow` bigint(20) NOT NULL DEFAULT '0',
  `pre_workflow` bigint(20) NOT NULL DEFAULT '0',
  `pre_workflow_info` varchar(255) DEFAULT NULL,
  `version` varchar(255) DEFAULT NULL,
  `version_code` varchar(255) DEFAULT NULL,
  `sequence` bigint(20) NOT NULL DEFAULT '0',
  `run_state` bigint(20) DEFAULT NULL,
  `event` bigint(20) DEFAULT '0',
  `manifest` longtext,
  `description` text,
  `source_info` longtext,
  `env` longtext,
  `requires` longtext,
  `auth_list` longtext,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `pre_stage` bigint(20) NOT NULL DEFAULT '0',
  `pre_action` bigint(20) NOT NULL DEFAULT '0',
  `fail_reason` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_workflow_log_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=621 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for workflow_sequence
-- ----------------------------
DROP TABLE IF EXISTS `workflow_sequence`;
CREATE TABLE `workflow_sequence` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `workflow` bigint(20) NOT NULL DEFAULT '0',
  `sequence` bigint(20) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_workflow_sequence_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=621 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for workflow_var
-- ----------------------------
DROP TABLE IF EXISTS `workflow_var`;
CREATE TABLE `workflow_var` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `workflow` bigint(20) DEFAULT NULL,
  `key` varchar(255) DEFAULT NULL,
  `default` longtext,
  `vaule` longtext,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_workflow_var_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for workflow_var_log
-- ----------------------------
DROP TABLE IF EXISTS `workflow_var_log`;
CREATE TABLE `workflow_var_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `workflow` bigint(20) DEFAULT NULL,
  `from_workflow` bigint(20) DEFAULT NULL,
  `sequence` bigint(20) DEFAULT NULL,
  `key` varchar(255) DEFAULT NULL,
  `default` longtext,
  `vaule` longtext,
  `change_log` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

```
