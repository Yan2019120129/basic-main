package sql

import (
	"basic/tools/utils"
)

const UserCommentTableName = "user_comment"
const UserCommentTableComment = "商城财务统计表"
const CreateUserComment = `CREATE TABLE ` + UserCommentTableName + ` (
id bigint NOT NULL COMMENT '评论ID',
user_id bigint NOT NULL COMMENT '用户ID',
product_id bigint NOT NULL COMMENT '商品ID',
star_rating varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '星级',
username varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '用户名',
user_avatar varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '用户头像',
comment varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '评论',
time int NOT NULL COMMENT '时间',
status tinyint NOT NULL COMMENT '状态：1新增，2追加，3回复',
admin_id int DEFAULT NULL,
PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + UserCommentTableComment + `';`

var BasicUserComment = &utils.InitTable{
	Name:        UserCommentTableName,
	Comment:     UserCommentTableComment,
	CreateTable: CreateUserComment,
}
