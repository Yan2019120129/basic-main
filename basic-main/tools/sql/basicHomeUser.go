package sql

import (
	"basic/tools/utils"
)

const BasicHomeUserTableName = "user"
const BasicHomeUserTableComment = "前台用户"
const CreateBasicHomeUser = `CREATE TABLE ` + BasicHomeUserTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	parent_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '上级ID',
	country_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '国家ID',
	username VARCHAR(191) NOT NULL DEFAULT '' COMMENT '用户名',
	nickname VARCHAR(50) NOT NULL DEFAULT '' COMMENT '昵称',
	email VARCHAR(50) NOT NULL DEFAULT '' COMMENT '邮箱',
	telephone VARCHAR(50) NOT NULL DEFAULT '' COMMENT '手机号码',
	avatar VARCHAR(255) NOT NULL DEFAULT '' COMMENT '头像',
	sex TINYINT NOT NULL DEFAULT -1 COMMENT '性别 -1未知 1男 2女',
	birthday INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '生日',
	password VARCHAR(255) NOT NULL DEFAULT '' COMMENT '密码',
	security_key VARCHAR(255) NOT NULL DEFAULT '' COMMENT '安全密钥',
	money DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '金额',
	freeze_money DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '冻结金额',
	type TINYINT NOT NULL DEFAULT 10 COMMENT '类型 -10客服坐席 -1虚拟用户 10真实用户 20会员用户',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -2删除 -1禁用 10启用',
	data VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据',
	ip4 INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'IP4地址',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
	UNIQUE KEY user_username (username)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + BasicHomeUserTableComment + `';`

const InsertBasicHomeUser = `INSERT INTO ` + BasicHomeUserTableName + `(admin_id, username, nickname, avatar, password, security_key, type) VALUES
(1, 'ceshi', 'ceshi', '', '14e1b600b1fd579f47433b88e8d85291', '14e1b600b1fd579f47433b88e8d85291', 10),
(1, 'adminOnline', 'admin', '/assets/images/avatar/0.jpeg', '', '', -10);`

var BasicHomeUser = &utils.InitTable{
	Name:        BasicHomeUserTableName,
	Comment:     BasicHomeUserTableComment,
	CreateTable: CreateBasicHomeUser,
	InsertTable: InsertBasicHomeUser,
}
