package sql

import (
	"basic/tools/utils"
)

const StoreTableName = "store"
const StoreTableComment = "商城商店表"
const CreateStore = `CREATE TABLE ` + StoreTableName + ` (
  id bigint NOT NULL COMMENT '商店ID',
  user_id bigint NOT NULL COMMENT '用户ID',
  admin_id bigint NOT NULL COMMENT '管理员ID',
  sales_volume int DEFAULT NULL COMMENT '销售额',
  visitor_count int DEFAULT NULL COMMENT '访客数',
  order_count int DEFAULT NULL COMMENT '订单数',
  yesterday_difference int DEFAULT NULL COMMENT '昨日差',
  rating int DEFAULT NULL COMMENT '评分',
  pending_payment varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '待付款',
  pending_shipment varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '待发货',
  pending_receipt varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '待收货',
  after_sales_service varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '待售后',
  pending_review varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '待评论',
  store_logo varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '店铺log',
  store_name varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '店铺名称',
  phone varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '电话',
  store_type varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '类型',
  keywords varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '关键词',
  status tinyint NOT NULL DEFAULT '10' COMMENT '状态： -2关闭， -1整改维护， 1在使用， 10启用',
  description varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '描述'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + StoreTableComment + `';`

var BasicStore = &utils.InitTable{
	Name:        StoreTableName,
	Comment:     StoreTableComment,
	CreateTable: CreateStore,
}
