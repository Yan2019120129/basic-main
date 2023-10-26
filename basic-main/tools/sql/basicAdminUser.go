package sql

import (
	"basic/tools/utils"
)

const AdminUserTableName = "admin_user"
const AdminUserTableComment = "后台管理"
const CreateAdminUser = `CREATE TABLE ` + AdminUserTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	parent_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '上级ID',
	username VARCHAR(50) NOT NULL DEFAULT '' COMMENT '用户名',
	email VARCHAR(50) NOT NULL DEFAULT '' COMMENT '邮件',
	nickname VARCHAR(50) NOT NULL DEFAULT '' COMMENT '昵称', 
	avatar VARCHAR(50) NOT NULL DEFAULT '' COMMENT '头像',
	password VARCHAR(255) NOT NULL DEFAULT '' COMMENT '密码',
	security_key VARCHAR(255) NOT NULL DEFAULT '' COMMENT '安全密钥',
	money DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '金额',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -2删除 -1禁用 10启用',
	data VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据',
	extra VARCHAR(255) NOT NULL DEFAULT '' COMMENT '额外',
	domain VARCHAR(255) NOT NULL DEFAULT '' COMMENT '域名',
	expired_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '过期时间',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
	UNIQUE KEY admin_user_username (username)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + AdminUserTableComment + `';`

const InsertAdminUser = `INSERT INTO admin_user(username, password, security_key, data) VALUES
('admin', '14e1b600b1fd579f47433b88e8d85291', '14e1b600b1fd579f47433b88e8d85291', '{"key": "8888", "only": true, "expire": 86400, "whitelist": "", "blacklist": ""}');`

var BasicAdminUser = &utils.InitTable{
	Name:        AdminUserTableName,
	Comment:     AdminUserTableComment,
	CreateTable: CreateAdminUser,
	InsertTable: InsertAdminUser,
}
