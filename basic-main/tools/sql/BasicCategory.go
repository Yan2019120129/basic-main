package sql

import (
	"basic/tools/utils"
)

const CategoryTableName = "category"
const CategoryTableComment = "商城类目表"
const CreateCategory = `CREATE TABLE ` + CategoryTableName + ` (
  id bigint NOT NULL COMMENT '类目ID',
  admin_id bigint NOT NULL COMMENT '管理员ID',
  image varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '类目图片',
  name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '种类名称',
  date int NOT NULL COMMENT '时间',
  status tinyint NOT NULL COMMENT '状态',
  operate tinyint NOT NULL COMMENT '操作'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + CategoryTableComment + `';`

var BasicCategory = &utils.InitTable{
	Name:        CategoryTableName,
	Comment:     CategoryTableComment,
	CreateTable: CreateCategory,
}
