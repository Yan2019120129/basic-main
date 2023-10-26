package sql

import (
	"basic/tools/utils"
)

const HomeUserAssetsTableName = "user_assets"
const HomeUserAssetsTableComment = "用户资产"
const CreateHomeUserAssets = `CREATE TABLE ` + HomeUserAssetsTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	user_id  INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
	assets_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '资产ID',
	money DECIMAL(20, 8) NOT NULL DEFAULT 0 COMMENT '金额',
	freeze_money DECIMAL(20, 8) NOT NULL DEFAULT 0 COMMENT '冻结金额',
	data VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -2删除｜-1禁用｜10启用',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + HomeUserAssetsTableComment + `';`

const InsertHomeUserAssets = ``

var BasicHomeUserAssets = &utils.InitTable{
	Name:        HomeUserAssetsTableName,
	Comment:     HomeUserAssetsTableComment,
	CreateTable: CreateHomeUserAssets,
	InsertTable: InsertHomeUserAssets,
}
