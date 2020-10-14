/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 80019
 Source Host           : localhost
 Source Database       : mp

 Target Server Type    : MySQL
 Target Server Version : 80019
 File Encoding         : utf-8

 Date: 10/14/2020 10:20:52 AM
*/

SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `privilege_admin`
-- ----------------------------
DROP TABLE IF EXISTS `privilege_admin`;
CREATE TABLE `privilege_admin` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '后台账户ID',
  `login_name` varchar(32) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT '登录名',
  `password` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT '密码hash',
  `username` varchar(128) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT '用户名',
  `avatar` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '管理员图像',
  `salt` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT '盐值',
  `dep_key` varchar(40) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT '部门key(3位一个)',
  `role_key` varchar(20) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT '角色key(3位一个)',
  `is_admin` tinyint NOT NULL DEFAULT '1' COMMENT '是否是超管 1:超管 2:默认非超管',
  `is_delete` tinyint NOT NULL DEFAULT '2' COMMENT '状态 1:删除 2:正常',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '用户创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='后台管理员';

-- ----------------------------
--  Records of `privilege_admin`
-- ----------------------------
BEGIN;
INSERT INTO `privilege_admin` VALUES ('1', 'admin', 'f44d47497cf85656fedf5ef3c12ed7e1', 'admin', '', '1515317636278', '', '000', '0', '2', '2020-05-08 13:56:35', '2020-06-16 20:24:58'), ('2', 'admin1', 'f44d47497cf85656fedf5ef3c12ed7e1', 'admin1', '', '1515317636278', '', '103', '0', '2', '2020-05-08 16:15:06', '2020-06-16 20:24:58'), ('3', 'dyong', '3e96f0307b71c864b3de289d40feeafe', 'dyong@123', '', '1519768484413', '', '100', '0', '2', '2020-05-11 19:01:09', '2020-07-01 16:39:18'), ('6', 'hzai', '36bff3582c042ba29e21a7f42ce05c2d', 'huangzai', '', '1516389325208', '', '', '0', '2', '2020-05-21 16:24:13', '2020-06-16 20:24:58'), ('9', 'zhangsan', '6884dd22b0aa0d070c86da24be68356d', '李四', '', '1516402476669', '', '101', '0', '2', '2020-05-21 20:13:19', '2020-06-16 20:24:58'), ('10', 'cs1', 'f8d173afa2c1c8fe72c60ab868e32edb', '田秋艳', '', '1519676431185', '', '101', '0', '2', '2020-06-30 13:50:34', '2020-06-30 13:50:34'), ('11', 'cs2', 'af9996fb1ccabfb379e5b9b0850dcb66', '田秋艳02', '', '1519676684269', '', '103', '0', '2', '2020-06-30 13:54:59', '2020-06-30 13:54:59'), ('13', 'sunchao', '70910d31f891b68f07541c0956ce08c4', '孙超', '', '1519686635019', '', '100', '0', '2', '2020-06-30 16:48:53', '2020-06-30 16:48:53'), ('15', 'lyuan', 'eff9c12a055d29e784e3af38e61142ac', '廖愿', '', '1519853351731', '', '107', '0', '2', '2020-07-02 17:22:28', '2020-07-02 18:31:06'), ('16', 'tina', '649a37ed3ef91bed493903ed14229eb6', 'tina', '', '1519912743857', '', '107', '0', '2', '2020-07-03 10:40:25', '2020-07-03 10:40:25');
COMMIT;

-- ----------------------------
--  Table structure for `privilege_menu`
-- ----------------------------
DROP TABLE IF EXISTS `privilege_menu`;
CREATE TABLE `privilege_menu` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `path` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '菜单前台路由',
  `component` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '组件唯一标识',
  `title` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '菜单名称',
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '前端vue组件名称',
  `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '图标样式',
  `menu_key` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '菜单key，三位一层，最多支持5级菜单',
  `level` tinyint NOT NULL DEFAULT '0' COMMENT '菜单层级',
  `display_order` int NOT NULL DEFAULT '1' COMMENT '菜单排序',
  `is_hidden` int NOT NULL DEFAULT '2' COMMENT '是否显示该菜单项 1:隐藏 2:显示 ',
  `is_delete` tinyint NOT NULL DEFAULT '2' COMMENT '状态 1:删除 2:正常',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `udx_key` (`menu_key`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='菜单';

-- ----------------------------
--  Records of `privilege_menu`
-- ----------------------------
BEGIN;
INSERT INTO `privilege_menu` VALUES ('1', '/course', 'Layout', '课程管理', 'Course', 'course', '100', '1', '3', '2', '1', '2020-05-08 17:48:04', '2020-10-13 15:00:54'), ('2', 'course-management', 'course/management/coursemanagement', '章节管理', 'CourseManagement', 'course-module', '100100', '2', '1', '2', '1', '2020-05-08 17:48:53', '2020-10-13 15:00:54'), ('3', 'lesson-package', 'course/package/coursepackage', '课包管理', 'LessonPackage', 'course-package', '100200', '2', '2', '2', '1', '2020-05-09 13:11:47', '2020-10-13 15:00:54'), ('4', '/sys', 'Layout', '管理员管理', 'Sys', 'admin', '200', '1', '2', '2', '2', '2020-05-09 13:12:14', '2020-06-16 21:20:53'), ('5', 'sys-list', 'sys/list', '管理员列表', 'SysList', 'list', '200100', '2', '1', '2', '2', '2020-05-09 13:12:39', '2020-06-16 21:20:56'), ('6', 'sys-role', 'sys/role', '角色管理', 'SysRoles', 'roles', '200200', '2', '2', '2', '2', '2020-05-09 13:13:11', '2020-06-16 21:20:58'), ('9', '/user', 'Layout', '用户管理', 'User', 'user', '300', '1', '1', '2', '2', '2020-05-09 20:49:41', '2020-06-16 21:21:00'), ('11', 'course-chapter', 'course/chapter/coursechapter', '章节管理', 'CourseChapter', 'course-package', '100300', '2', '3', '2', '1', '2020-05-18 14:31:04', '2020-10-13 15:00:54'), ('12', 'management', 'user/userlist', '用户管理', 'UserList', 'user', '300100', '2', '2', '2', '2', '2020-05-18 18:40:52', '2020-06-16 21:21:05'), ('13', 'content-management', 'course/content/contentmanagement', '内容管理', 'ContentManagement', 'course-package', '100400', '2', '4', '2', '1', '2020-05-21 20:36:10', '2020-10-13 15:00:54'), ('25', 'sys-menus', 'sys/menu', '菜单管理', 'SysMenus', 'list', '200300', '2', '3', '2', '2', '2020-05-25 17:01:15', '2020-06-16 21:21:11'), ('26', '/basic', 'Layout', '基础管理', 'Basic', 'basic', '300101', '2', '4', '2', '2', '2020-05-27 18:19:44', '2020-06-16 21:21:13');
COMMIT;

-- ----------------------------
--  Table structure for `privilege_menu_back_url`
-- ----------------------------
DROP TABLE IF EXISTS `privilege_menu_back_url`;
CREATE TABLE `privilege_menu_back_url` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '路由表',
  `menu_key` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '菜单menu_key',
  `back_url` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '后台路由',
  `desc` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '描述',
  `is_delete` tinyint NOT NULL DEFAULT '2' COMMENT '状态 1:删除 2:正常',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `udx_url` (`menu_key`,`back_url`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='前端菜单与后台路由映射关系';

-- ----------------------------
--  Records of `privilege_menu_back_url`
-- ----------------------------
BEGIN;
INSERT INTO `privilege_menu_back_url` VALUES ('1', '100', '/privilege/menu/list', '', '2', '2020-05-09 15:39:55', '2020-06-16 21:18:11'), ('3', '100100', '/privilege/menu/add', '', '2', '2020-05-09 15:40:45', '2020-06-16 21:18:11'), ('4', '200', '/privilege/menu/update', '管理员管理更新', '2', '2020-05-09 15:40:54', '2020-07-10 11:37:45'), ('5', '200', '/privilege/menu/list/tree', '管理员菜单列表', '2', '2020-05-09 15:41:07', '2020-07-10 11:38:11');
COMMIT;

-- ----------------------------
--  Table structure for `privilege_role`
-- ----------------------------
DROP TABLE IF EXISTS `privilege_role`;
CREATE TABLE `privilege_role` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '角色名称',
  `desc` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '描述',
  `role_key` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '角色key',
  `status` tinyint NOT NULL DEFAULT '2' COMMENT '状态 1:未启用 2:已启用',
  `is_delete` tinyint NOT NULL DEFAULT '2' COMMENT '状态 1:删除 2:正常',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `udx_key_name` (`role_key`,`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
--  Records of `privilege_role`
-- ----------------------------
BEGIN;
INSERT INTO `privilege_role` VALUES ('19', '角色管理员', '管理角色', '102', '2', '2', '2020-05-21 16:21:35', '2020-06-16 21:18:18'), ('21', '课程管理员', '管理课程', '103', '2', '2', '2020-05-21 16:41:21', '2020-06-16 21:18:18'), ('34', '课包管理员', '管理课包', '104', '2', '2', '2020-05-22 13:37:05', '2020-06-16 21:18:18'), ('35', '内容管理员', '管理内容', '105', '2', '2', '2020-05-26 10:27:36', '2020-06-16 21:18:18'), ('36', '角色测试01-管理员管理', '角色测试01', '106', '2', '0', '2020-07-02 17:21:51', '2020-07-02 17:21:51'), ('37', '全部模快', '全部模快', '107', '2', '0', '2020-07-02 18:30:55', '2020-07-02 18:30:55');
COMMIT;

-- ----------------------------
--  Table structure for `privilege_role_menu`
-- ----------------------------
DROP TABLE IF EXISTS `privilege_role_menu`;
CREATE TABLE `privilege_role_menu` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `role_key` varchar(15) NOT NULL DEFAULT '' COMMENT '角色key',
  `menu_key` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '菜单menu_key',
  `is_delete` tinyint NOT NULL DEFAULT '2' COMMENT '状态 1:删除 2:正常',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_key` (`role_key`),
  KEY `idx_mkey` (`menu_key`)
) ENGINE=InnoDB AUTO_INCREMENT=163 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色菜单关联关系';

-- ----------------------------
--  Records of `privilege_role_menu`
-- ----------------------------
BEGIN;
INSERT INTO `privilege_role_menu` VALUES ('7', '200', '200', '2', '2020-05-11 14:14:52', '2020-06-16 21:18:23'), ('8', '200', '200100', '2', '2020-05-15 00:10:20', '2020-06-16 21:18:23'), ('9', '200', '200200', '2', '2020-05-15 00:10:25', '2020-06-16 21:18:23'), ('10', '200', '300', '2', '2020-05-15 00:10:42', '2020-06-16 21:18:23'), ('11', '200', '300100', '2', '2020-05-19 09:22:07', '2020-06-16 21:18:23'), ('12', '200003003', '200200', '2', '2020-05-21 15:20:33', '2020-06-16 21:18:23'), ('13', '200003003', '200100', '2', '2020-05-21 15:20:33', '2020-06-16 21:18:23'), ('14', '200003003', '200', '2', '2020-05-21 15:20:33', '2020-06-16 21:18:23'), ('41', '100', '300100', '2', '2020-05-21 16:27:27', '2020-06-16 21:18:23'), ('42', '100', '100300', '2', '2020-05-21 16:27:27', '2020-06-16 21:18:23'), ('43', '100', '300', '2', '2020-05-21 16:27:27', '2020-06-16 21:18:23'), ('44', '100', '100', '2', '2020-05-21 16:27:27', '2020-06-16 21:18:23'), ('45', '100', '200100', '2', '2020-05-21 16:27:27', '2020-06-16 21:18:23'), ('46', '100', '200200', '2', '2020-05-21 16:27:27', '2020-06-16 21:18:23'), ('47', '100', '100200', '2', '2020-05-21 16:27:27', '2020-06-16 21:18:23'), ('48', '100', '100100', '2', '2020-05-21 16:27:27', '2020-06-16 21:18:23'), ('49', '100', '200', '2', '2020-05-21 16:27:27', '2020-06-16 21:18:23'), ('56', '101', '300', '2', '2020-05-21 16:57:36', '2020-06-16 21:18:23'), ('57', '101', '300100', '2', '2020-05-21 16:57:36', '2020-06-16 21:18:23'), ('63', '100', '100400', '2', '2020-05-21 20:39:27', '2020-06-16 21:18:23'), ('64', '200', '100400', '2', '2020-05-21 20:39:43', '2020-06-16 21:18:23'), ('102', '102', '200', '2', '2020-05-22 13:33:25', '2020-06-16 21:18:23'), ('103', '102', '200200', '2', '2020-05-22 13:33:25', '2020-06-16 21:18:23'), ('108', '104', '100', '2', '2020-05-22 13:37:05', '2020-06-16 21:18:23'), ('109', '104', '100200', '2', '2020-05-22 13:37:05', '2020-06-16 21:18:23'), ('130', '105', '100', '2', '2020-05-26 10:27:36', '2020-06-16 21:18:23'), ('131', '105', '100400', '2', '2020-05-26 10:27:36', '2020-06-16 21:18:23'), ('132', '100', '200300', '2', '2020-05-26 11:40:38', '2020-06-16 21:18:23'), ('133', '200', '200300', '2', '2020-05-26 11:40:47', '2020-06-16 21:18:23'), ('134', '100', '300101', '2', '2020-05-27 18:24:30', '2020-06-16 21:18:23'), ('135', '103', '300', '2', '2020-06-02 16:17:26', '2020-06-16 21:18:23'), ('136', '103', '300101', '2', '2020-06-02 16:17:26', '2020-06-16 21:18:23'), ('137', '103', '100', '2', '2020-06-02 16:17:26', '2020-06-16 21:18:23'), ('138', '103', '200100', '2', '2020-06-02 16:17:26', '2020-06-16 21:18:23'), ('139', '103', '200', '2', '2020-06-02 16:17:26', '2020-06-16 21:18:23'), ('140', '103', '100400', '2', '2020-06-02 16:17:26', '2020-06-16 21:18:23'), ('141', '103', '200200', '2', '2020-06-02 16:17:26', '2020-06-16 21:18:23'), ('142', '103', '100100', '2', '2020-06-02 16:17:26', '2020-06-16 21:18:23'), ('143', '103', '300100', '2', '2020-06-02 16:17:26', '2020-06-16 21:18:23'), ('144', '103', '100200', '2', '2020-06-02 16:17:26', '2020-06-16 21:18:23'), ('145', '103', '200300', '2', '2020-06-02 16:17:26', '2020-06-16 21:18:23'), ('146', '103', '100300', '2', '2020-06-02 16:17:26', '2020-06-16 21:18:23'), ('147', '106', '200', '2', '2020-07-02 17:21:51', '2020-07-02 17:21:51'), ('148', '106', '200300', '2', '2020-07-02 17:21:51', '2020-07-02 17:21:51'), ('149', '106', '200100', '2', '2020-07-02 17:21:51', '2020-07-02 17:21:51'), ('150', '106', '200200', '2', '2020-07-02 17:21:51', '2020-07-02 17:21:51'), ('151', '107', '300100', '2', '2020-07-02 18:30:55', '2020-07-02 18:30:55'), ('152', '107', '300', '2', '2020-07-02 18:30:55', '2020-07-02 18:30:55'), ('153', '107', '100400', '2', '2020-07-02 18:30:55', '2020-07-02 18:30:55'), ('154', '107', '300101', '2', '2020-07-02 18:30:55', '2020-07-02 18:30:55'), ('155', '107', '100100', '2', '2020-07-02 18:30:55', '2020-07-02 18:30:55'), ('156', '107', '200200', '2', '2020-07-02 18:30:55', '2020-07-02 18:30:55'), ('157', '107', '200100', '2', '2020-07-02 18:30:55', '2020-07-02 18:30:55'), ('158', '107', '200300', '2', '2020-07-02 18:30:55', '2020-07-02 18:30:55'), ('159', '107', '200', '2', '2020-07-02 18:30:55', '2020-07-02 18:30:55'), ('160', '107', '100300', '2', '2020-07-02 18:30:55', '2020-07-02 18:30:55'), ('161', '107', '100', '2', '2020-07-02 18:30:55', '2020-07-02 18:30:55'), ('162', '107', '100200', '2', '2020-07-02 18:30:55', '2020-07-02 18:30:55');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
