package sql

import (
	"basic/tools/utils"
)

const CommodityTableName = "commodity"
const CommodityTableComment = "商城商品表"
const CreateCommodity = `CREATE TABLE ` + CommodityTableName + ` (
  id int NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  store_id bigint NOT NULL COMMENT '店铺ID',
  admin_id bigint NOT NULL COMMENT '管理员ID',
  product_image varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '商品图片',
  name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '名称',
  purchase_price decimal(10,2) NOT NULL COMMENT '进货价',
  selling_price decimal(10,2) NOT NULL COMMENT '出售价',
  stock int NOT NULL COMMENT '库存',
  sales_volume int NOT NULL COMMENT '出货量',
  status tinyint NOT NULL COMMENT '状态 -1删除 -2下架 1在售',
  operation tinyint NOT NULL COMMENT '操作',
  category_id bigint NOT NULL COMMENT '类目ID',
  commodity_id bigint NOT NULL COMMENT '分类ID',
  specification_attributes varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '规格属性',
  description varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '描述',
  brand varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '品牌',
  time int NOT NULL COMMENT '时间',
  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + CommodityTableComment + `';`

var BasicCommodity = &utils.InitTable{
	Name:        CommodityTableName,
	Comment:     CommodityTableComment,
	CreateTable: CreateCommodity,
}
