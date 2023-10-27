package sql

import (
	"basic/tools/utils"
)

const ShopMyWalletTableName = "shop_my_wallet"
const ShopMyWalletTableComment = "商城用户钱包记录"
const CreateShopMyWallet = `CREATE TABLE ` + ShopMyWalletTableName + ` (
  id bigint NOT NULL COMMENT '钱包唯一ID',
  user_id bigint NOT NULL COMMENT '用户ID',
  amount decimal(10,0) NOT NULL COMMENT '金额',
  payment_method varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '支付方式',
  date int DEFAULT NULL COMMENT '交易时间',
  status tinyint DEFAULT NULL COMMENT '支付状态：-1待支付，1已支付',
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + ShopMyWalletTableComment + `';`

var BasicShopMyWallet = &utils.InitTable{
	Name:        ShopMyWalletTableName,
	Comment:     ShopMyWalletTableComment,
	CreateTable: CreateShopMyWallet,
}
