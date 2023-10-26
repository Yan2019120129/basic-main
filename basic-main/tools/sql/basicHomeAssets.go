package sql

import (
	"basic/tools/utils"
)

const HomeAssetsTableName = "assets"
const HomeAssetsTableComment = "平台资产"
const CreateHomeAssets = `CREATE TABLE ` + HomeAssetsTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	name VARCHAR(191) NOT NULL DEFAULT '' COMMENT '名称',
	icon VARCHAR(191) NOT NULL DEFAULT '' COMMENT '图标',
	type TINYINT NOT NULL DEFAULT 1 COMMENT '类型 1数字金额 10银行金额 20奖励金额',
	data VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -2删除｜-1禁用｜10启用',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + HomeAssetsTableComment + `';`

const InsertHomeAssets = `INSERT INTO ` + HomeAssetsTableName + `(admin_id, name, icon, type, status) VALUES
(1, 'digitalAmount', '/assets/images/icon/currency/usdt.png', 1, 10),
(1, 'banksAmount', '/assets/images/icon/bank_cards.png', 10, -1),
(1, 'rewardsAmount', '/assets/images/icon/award.png', 20, -1);`

var BasicHomeAssets = &utils.InitTable{
	Name:        HomeAssetsTableName,
	Comment:     HomeAssetsTableComment,
	CreateTable: CreateHomeAssets,
	InsertTable: InsertHomeAssets,
}
