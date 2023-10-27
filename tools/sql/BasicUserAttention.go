package sql

import (
	"basic/tools/utils"
)

const UserAttentionTableName = "user_attention"
const UserAttentionTableComment = "用户关注表"
const CreateUserAttention = `CREATE TABLE ` + UserAttentionTableName + ` (
  id bigint NOT NULL COMMENT '关注ID',
  user_id bigint NOT NULL COMMENT '用户ID',
  admin_id bigint NOT NULL COMMENT '管理员ID',
  shop_id bigint NOT NULL COMMENT '店铺ID',
  shop_logo varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '店铺log',
  shop_name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '店铺名',
  date int NOT NULL COMMENT '时间',
  status tinyint NOT NULL DEFAULT '1' COMMENT '状态：1关注，2取消关注'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + UserAttentionTableComment + `';`

var BasicUserAttention = &utils.InitTable{
	Name:        UserAttentionTableName,
	Comment:     UserAttentionTableComment,
	CreateTable: CreateUserAttention,
}
