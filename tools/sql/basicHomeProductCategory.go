package sql

import (
	"basic/tools/utils"
)

const HomeProductCategoryTableName = "product_category"
const HomeProductCategoryTableComment = "产品分类"
const CreateHomeProductCategory = `CREATE TABLE ` + HomeProductCategoryTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	parent_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '分类上级ID',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	type TINYINT NOT NULL DEFAULT 1 COMMENT '类型 1普通类型',
	name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '标题',
	image VARCHAR(255) NOT NULL DEFAULT '' COMMENT '封面',
	sort SMALLINT UNSIGNED NOT NULL DEFAULT 99 COMMENT '排序',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -2删除 -1禁用 10启用',
	recommend TINYINT NOT NULL DEFAULT -1 COMMENT '推荐 -1关闭 10推荐',
	data VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + HomeProductCategoryTableComment + `';`

const InsertHomeProductCategory = ``

var BasicHomeProductCategory = &utils.InitTable{
	Name:        HomeProductCategoryTableName,
	Comment:     HomeProductCategoryTableComment,
	CreateTable: CreateHomeProductCategory,
	InsertTable: InsertHomeProductCategory,
}
