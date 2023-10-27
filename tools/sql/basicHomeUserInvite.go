package sql

import (
	"basic/tools/utils"
)

const BasicHomeUserInviteTableName = "user_invite"
const BasicHomeUserInviteTableComment = "用户邀请"
const CreateBasicHomeUserInvite = `CREATE TABLE ` + BasicHomeUserInviteTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	user_id  INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
	code VARCHAR(64) NOT NULL DEFAULT '' COMMENT '邀请码',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -1禁用 10启用',
	data VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
	UNIQUE KEY user_invite_code (code)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + BasicHomeUserInviteTableComment + `';`

const InsertBasicHomeUserInvite = `INSERT INTO user_invite(admin_id, code) VALUES (1, '8888')`

var BasicHomeUserInvite = &utils.InitTable{
	Name:        BasicHomeUserInviteTableName,
	Comment:     BasicHomeUserInviteTableComment,
	CreateTable: CreateBasicHomeUserInvite,
	InsertTable: InsertBasicHomeUserInvite,
}
