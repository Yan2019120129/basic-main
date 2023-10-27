package sql

import (
	"basic/tools/utils"
)

const BasicHomeUserBillTableName = "user_bill"
const BasicHomeUserBillTableComment = "用户账单"
const CreateBasicHomeUserBill = `CREATE TABLE ` + BasicHomeUserBillTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理ID',
	user_id  INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
	source_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '来源ID',
	name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '标题',
	type TINYINT NOT NULL DEFAULT 0 COMMENT '类型',
	balance DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '余额',
	money DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '金额',
	data VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + BasicHomeUserBillTableComment + `';`

const InsertBasicHomeUserBill = ``

var BasicHomeUserBill = &utils.InitTable{
	Name:        BasicHomeUserBillTableName,
	Comment:     BasicHomeUserBillTableComment,
	CreateTable: CreateBasicHomeUserBill,
	InsertTable: InsertBasicHomeUserBill,
}
