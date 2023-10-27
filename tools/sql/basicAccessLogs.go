package sql

import (
	"basic/tools/utils"
)

const AccessLogsTableName = "access_logs"
const AccessLogsTableComment = "访问日志"
const CreateAccessLogs = `CREATE TABLE ` + AccessLogsTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	user_id  INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
	type TINYINT NOT NULL DEFAULT 1 COMMENT '类型 1后端日志 2前端日志',
	name VARCHAR(64) NOT NULL DEFAULT '' COMMENT '标题',
	ip4 INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'IP4地址',
	user_agent VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'ua信息',
	lang VARCHAR(255) NOT NULL DEFAULT '' COMMENT '语言信息',
	route VARCHAR(64) NOT NULL DEFAULT '' COMMENT '操作路由',
	data TEXT COMMENT '数据',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + AccessLogsTableComment + `';`

const InsertAccessLogs = ``

var BasicAccessLogs = &utils.InitTable{
	Name:        AccessLogsTableName,
	Comment:     AccessLogsTableComment,
	CreateTable: CreateAccessLogs,
	InsertTable: InsertAccessLogs,
}
