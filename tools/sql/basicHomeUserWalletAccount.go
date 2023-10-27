package sql

import (
	"basic/tools/utils"
)

const HomeUserWalletAccountTableName = "user_wallet_account"
const HomeUserWalletAccountTableComment = "钱包账号"
const CreateHomeUserWalletAccount = `CREATE TABLE ` + HomeUserWalletAccountTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	user_id  INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
	payment_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '提现方式ID',
	name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '名称',
	real_name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '真实名字',
	card_number VARCHAR(255) NOT NULL DEFAULT '' COMMENT '卡号',
	address VARCHAR(255) NOT NULL DEFAULT '' COMMENT '地址',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -2删除 -1禁用 10启用',
	sort TINYINT NOT NULL DEFAULT 0 COMMENT '排序',
	data VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + HomeUserWalletAccountTableComment + `';`

const InsertHomeUserWalletAccount = ``

var BasicHomeUserWalletAccount = &utils.InitTable{
	Name:        HomeUserWalletAccountTableName,
	Comment:     HomeUserWalletAccountTableComment,
	CreateTable: CreateHomeUserWalletAccount,
	InsertTable: InsertHomeUserWalletAccount,
}
