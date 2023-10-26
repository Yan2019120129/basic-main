package sql

import (
	"basic/tools/utils"
)

const HomeProductOrderTableName = "product_order"
const HomeProductOrderTableComment = "产品订单"
const CreateHomeProductOrder = `CREATE TABLE ` + HomeProductOrderTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	user_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
	product_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '产品ID',
	order_sn VARCHAR(255) NOT NULL DEFAULT '' COMMENT '订单编号',
	money DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '金额',
	nums SMALLINT UNSIGNED NOT NULL DEFAULT 99 COMMENT '数量',
	type TINYINT NOT NULL DEFAULT 1 COMMENT '类型',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -2删除 -1完结 10启用',
	data VARCHAR(2048) NOT NULL DEFAULT '' COMMENT '数据',
	expired_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '过期时间',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + HomeProductOrderTableComment + `';`

const InsertHomeProductOrder = ``

var BasicHomeProductOrder = &utils.InitTable{
	Name:        HomeProductOrderTableName,
	Comment:     HomeProductOrderTableComment,
	CreateTable: CreateHomeProductOrder,
	InsertTable: InsertHomeProductOrder,
}
