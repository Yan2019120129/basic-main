package sql

import (
	"basic/tools/utils"
)

const CommodityOrderTableName = "commodity_order"
const CommodityOrderTableComment = "商城商品订单信息表"
const CreateCommodityOrder = `CREATE TABLE ` + CommodityOrderTableName + ` (
  id bigint NOT NULL COMMENT '订单ID',
  admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
  shop_id bigint NOT NULL COMMENT '店铺ID',
  shop_logo varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '店铺LOG',
  shop_name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '店铺名',
  product_id bigint NOT NULL COMMENT '商品ID',
  product_image varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '商品图片',
  product_description varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '商品描述',
  attributes_specifications varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '属性价格',
  original_price decimal(10,2) NOT NULL COMMENT '原价',
  transaction_price decimal(10,2) NOT NULL COMMENT '成交价',
  quantity bigint NOT NULL COMMENT '数量',
  user_id bigint NOT NULL COMMENT '用户ID',
  payment_method varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '支付方式',
  shipping_address varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '收货地址',
  user_operate tinyint NOT NULL COMMENT '用户操作：1已下单，2取消订单',
  order_status tinyint NOT NULL COMMENT '订单状态：0待处理，1已支付，2已发货，3已送达，-1已取消',
  user_operate_time int NOT NULL COMMENT '用户操作时间',
  store_id bigint NOT NULL COMMENT '店家ID',
  store_operate tinyint NOT NULL DEFAULT '1' COMMENT '店家操作：1接受订单操作，2拒绝订单操作，3发货操作，4取消订单操作',
  store_operate_time int NOT NULL COMMENT '店家操作时间'
  admin_operate_time int NOT NULL COMMENT '管理员操作时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + CommodityOrderTableComment + `';`

var BasicCommodityOrder = &utils.InitTable{
	Name:        CommodityOrderTableName,
	Comment:     CommodityOrderTableComment,
	CreateTable: CreateCommodityOrder,
}
