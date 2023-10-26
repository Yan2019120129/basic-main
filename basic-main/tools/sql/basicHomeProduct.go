package sql

import (
	"basic/tools/utils"
)

const HomeProductTableName = "product"
const HomeProductTableComment = "产品商品"
const CreateHomeProduct = `CREATE TABLE ` + HomeProductTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	category_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '类目ID',
	assets_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '资产ID',
	name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '标题',
	images VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '图片列表',
	money DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '金额',
	type TINYINT NOT NULL DEFAULT 1 COMMENT '类型 1普通类型',
	sort SMALLINT UNSIGNED NOT NULL DEFAULT 99 COMMENT '排序',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -2删除 -1禁用 10启用',
	recommend TINYINT NOT NULL DEFAULT -1 COMMENT '推荐 -1关闭 10推荐',
	sales INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '销售量',
	nums TINYINT NOT NULL DEFAULT -1 COMMENT '限购 -1 无限制',
	used INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '已用',
	total INT UNSIGNED NOT NULL DEFAULT 1000 COMMENT '总数',
	data VARCHAR(2048) NOT NULL DEFAULT '' COMMENT '数据',
	describes VARCHAR(255) NOT NULL DEFAULT '' COMMENT '描述',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + HomeProductTableComment + `';`

const InsertHomeProduct = ``

var BasicHomeProduct = &utils.InitTable{
	Name:        HomeProductTableName,
	Comment:     HomeProductTableComment,
	CreateTable: CreateHomeProduct,
	InsertTable: InsertHomeProduct,
}
