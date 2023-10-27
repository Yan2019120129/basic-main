package sql

import (
	"basic/tools/utils"
)

const UserAddressTableName = "user_address"
const UserAddressTableComment = "收货地址表"
const CreateUserAddress = `CREATE TABLE ` + UserAddressTableName + ` (
  id int NOT NULL COMMENT '收货地址ID',
  user_id int NOT NULL COMMENT '用户ID',
  admin_id bigint NOT NULL COMMENT '管理员ID',
  name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '收货人名',
  phone varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '电话',
  country varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '国家',
  shipping_address varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '收货地址',
  door_number int NOT NULL COMMENT '门牌号',
  zip_code int NOT NULL COMMENT '邮编',
  is_default tinyint NOT NULL COMMENT '是否默认：1是，-1否',
  time int NOT NULL COMMENT '时间',
  status tinyint NOT NULL COMMENT '状态：1在使用，-1已删除',
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + UserAddressTableComment + `';`

var BasicUserAddress = &utils.InitTable{
	Name:        UserAddressTableName,
	Comment:     UserAddressTableComment,
	CreateTable: CreateUserAddress,
}
