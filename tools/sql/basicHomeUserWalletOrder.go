package sql

import (
	"basic/tools/utils"
)

const HomeUserWalletOrderTableName = "user_wallet_order"
const HomeUserWalletOrderTableComment = "钱包订单"
const CreateHomeUserWalletOrder = `CREATE TABLE ` + HomeUserWalletOrderTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	order_sn VARCHAR(64) NOT NULL DEFAULT '' COMMENT '订单号',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	user_id  INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
	user_type TINYINT NOT NULL DEFAULT 0 COMMENT '用户类型',
	type TINYINT NOT NULL DEFAULT 0 COMMENT '类型 1充值 2系统加款 10提现 11系统减款',
	payment_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '充值类型ID｜提现账户ID',
	money DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '金额',
	balance DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '余额',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -2删除 -1拒绝 10处理 20完成',
	proof VARCHAR(255) NOT NULL DEFAULT '' COMMENT '凭证',
	data VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据',
	fee DECIMAL(6, 2) NOT NULL DEFAULT 0 COMMENT '手续费',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
	UNIQUE KEY user_wallet_order_sn (order_sn)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + HomeUserWalletOrderTableComment + `';`

const InsertHomeUserWalletOrder = ``

var BasicHomeUserWalletOrder = &utils.InitTable{
	Name:        HomeUserWalletOrderTableName,
	Comment:     HomeUserWalletOrderTableComment,
	CreateTable: CreateHomeUserWalletOrder,
	InsertTable: InsertHomeUserWalletOrder,
}
