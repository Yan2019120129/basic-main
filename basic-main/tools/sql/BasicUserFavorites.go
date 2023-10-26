package sql

import (
	"basic/tools/utils"
)

const UserFavoritesTableName = "user_favorites"
const UserFavoritesTableComment = "用户收藏表"
const CreateUserFavorites = `CREATE TABLE ` + UserFavoritesTableName + ` (
  id bigint NOT NULL COMMENT '收藏ID',
  user_id bigint NOT NULL COMMENT '用户ID',
  admin_id bigint NOT NULL COMMENT '管理员ID',
  commodity_id bigint NOT NULL COMMENT '商品ID',
  product_name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '商品名称',
  product_image varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '商品名称',
  date int NOT NULL COMMENT '收藏时间',
  status tinyint NOT NULL DEFAULT '1' COMMENT '收藏状态：-1取消收藏，1收藏',
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + UserFavoritesTableComment + `';`

var BasicUserFavorites = &utils.InitTable{
	Name:        UserFavoritesTableName,
	Comment:     UserFavoritesTableComment,
	CreateTable: CreateUserFavorites,
}
