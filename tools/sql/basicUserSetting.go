package sql

import (
	"basic/tools/utils"
)

const UserSettingTableName = "user_setting"
const UserSettingTableComment = "用户设置"
const CreateUserSetting = `CREATE TABLE ` + UserSettingTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理ID',
	user_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
	group_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '组ID',
	name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '名称',
	type VARCHAR(64) NOT NULL DEFAULT '' COMMENT '类型',
	field VARCHAR(64) NOT NULL DEFAULT '' COMMENT '健名',
	value VARCHAR(3072) NOT NULL DEFAULT '' COMMENT '健值',
	data VARCHAR(3072) NOT NULL DEFAULT '' COMMENT '数据',
	KEY admin_setting_key (field)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + UserSettingTableComment + `';`

const InsertUserSetting = ``

var BasicUserSetting = &utils.InitTable{
	Name:        UserSettingTableName,
	Comment:     UserSettingTableComment,
	CreateTable: CreateUserSetting,
	InsertTable: InsertUserSetting,
}
