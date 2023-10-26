package sql

import (
	"basic/tools/utils"
)

const ShopStatisticalRecordTableName = "shop_statistical_record"
const ShopStatisticalRecordTableComment = "商城店铺统计记录表"
const CreateShopStatisticalRecord = `CREATE TABLE ` + ShopStatisticalRecordTableName + ` (
  id bigint NOT NULL COMMENT '统计记录ID',
  shop_id bigint NOT NULL COMMENT '店铺ID',
  admin bigint NOT NULL COMMENT '管理员ID',
  visitor_count int DEFAULT NULL COMMENT '访客数',
  order_count int DEFAULT NULL COMMENT '订单数',
  earnings decimal(12,2) DEFAULT NULL COMMENT '收益',
  shop_favorites_count int DEFAULT NULL COMMENT '店铺收藏量',
  credit int DEFAULT NULL COMMENT '信用',
  product_favorites_count int DEFAULT NULL COMMENT '商品收藏量',
  product_count int DEFAULT NULL COMMENT '商品数量',
  time int DEFAULT NULL COMMENT '时间',
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + ShopStatisticalRecordTableComment + `';`

var BasicShopStatisticalRecord = &utils.InitTable{
	Name:        ShopStatisticalRecordTableName,
	Comment:     ShopStatisticalRecordTableComment,
	CreateTable: CreateShopStatisticalRecord,
}
