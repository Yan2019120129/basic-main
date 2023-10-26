package sql

import (
	"basic/tools/utils"
)

const ShopUserTableName = "shop_user"
const ShopUserTableComment = "商城用户表"
const CreateShopUser = `CREATE TABLE ` + ShopUserTableName + ` (
  id bigint NOT NULL COMMENT '用户唯一Id',
  head_portrait varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '用户头像',
  name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '用户姓名',
  sex tinyint NOT NULL DEFAULT '-1' COMMENT '用户性别：-1未知，1男，2女',
  birthday int DEFAULT NULL COMMENT '用户生日（存储时间戳）',
  email varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '用户邮箱',
  phone varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '用户电话',
  update_time int DEFAULT NULL COMMENT '用户信息最新修改时间（时间戳格式）',
  create_time int DEFAULT NULL COMMENT '用户注册时间(时间戳格式)',
  status tinyint NOT NULL DEFAULT '10' COMMENT '用户状态：-2删除，-1禁用，10启用',
  balance decimal(12,2) DEFAULT NULL COMMENT '用户余额',
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + ShopUserTableComment + `';`

var BasicShopUser = &utils.InitTable{
	Name:        ShopUserTableName,
	Comment:     ShopUserTableComment,
	CreateTable: CreateShopUser,
}
