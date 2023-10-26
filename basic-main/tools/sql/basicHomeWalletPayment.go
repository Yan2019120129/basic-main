package sql

import (
	"basic/tools/utils"
)

const HomeWalletPaymentTableName = "wallet_payment"
const HomeWalletPaymentTableComment = "支付方式"
const CreateHomeWalletPayment = `CREATE TABLE ` + HomeWalletPaymentTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	icon VARCHAR(255) NOT NULL DEFAULT '' COMMENT '图标',
	mode TINYINT NOT NULL DEFAULT 0 COMMENT '方式 1充值 10提现',
	type TINYINT NOT NULL DEFAULT 0 COMMENT '类型 1银行转账 10数字货币 20三方支付',
	name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '名称',
	account_name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '名字｜Token',
	account_code VARCHAR(255) NOT NULL DEFAULT '' COMMENT '卡号｜地址',
	sort TINYINT NOT NULL DEFAULT 0 COMMENT '排序',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -2删除 -1禁用 10启用',
	description VARCHAR(255) NOT NULL DEFAULT '' COMMENT '描述',
	data VARCHAR(3072) NOT NULL DEFAULT '' COMMENT '数据',
	expand VARCHAR(3072) NOT NULL DEFAULT '' COMMENT '扩展',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + HomeWalletPaymentTableComment + `';`

const InsertHomeWalletPayment = `INSERT INTO ` + HomeWalletPaymentTableName + `(admin_id, icon, mode, type, name, account_name, account_code, data, expand) VALUES
(1, '/assets/images/icon/bank_cards.png', 1, 1, 'Bank Card', 'Jon Lis', '88888888', '', ''),
(1, '/assets/images/icon/currency/usdt.png', 1, 10, 'TRC20', 'USDT', 'TEdAa5vLQ9pGqy3BB6NRDo1rQGvaeThpfF', '', ''),
(1, '/assets/images/icon/bank_cards.png', 10, 1, '银行卡', '', '', '', ''),
(1, '/assets/images/icon/currency/usdt.png', 10, 10, 'TRC20', 'USDT', '', '', '');`

var BasicHomeWalletPayment = &utils.InitTable{
	Name:        HomeWalletPaymentTableName,
	Comment:     HomeWalletPaymentTableComment,
	CreateTable: CreateHomeWalletPayment,
	InsertTable: InsertHomeWalletPayment,
}
