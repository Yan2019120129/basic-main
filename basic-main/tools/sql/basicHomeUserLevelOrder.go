package sql

import (
	"basic/tools/utils"
)

const HomeUserLevelOrderTableName = "user_level_order"
const HomeUserLevelOrderTableComment = "等级订单"
const CreateHomeUserLevelOrder = `CREATE TABLE ` + HomeUserLevelOrderTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	user_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
	level_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '等级ID',
	data VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '数据',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -1禁用 10启用',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + HomeUserLevelOrderTableComment + `';`

const InsertHomeUserLevelOrder = ``

var BasicHomeUserLevelOrder = &utils.InitTable{
	Name:        HomeUserLevelOrderTableName,
	Comment:     HomeUserLevelOrderTableComment,
	CreateTable: CreateHomeUserLevelOrder,
	InsertTable: InsertHomeUserLevelOrder,
}
