package sql

import (
	"basic/tools/utils"
)

const BasicChatConversationTableName = "chat_conversation"
const BasicChatConversationTableComment = "聊天会话"
const CreateBasicChatConversation = `CREATE TABLE ` + BasicChatConversationTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	conversation_id CHAR(32) NOT NULL COMMENT '会话ID',
	user_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
	receiver_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '接收者用户ID',
	type TINYINT NOT NULL DEFAULT 1 COMMENT '1在线客服',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -2删除 -1屏蔽 10正常',
	data VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
	KEY conversation_id (conversation_id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + BasicChatConversationTableComment + `';`

const InsertBasicChatConversation = ``

var BasicChatConversation = &utils.InitTable{
	Name:        BasicChatConversationTableName,
	Comment:     BasicChatConversationTableComment,
	CreateTable: CreateBasicChatConversation,
	InsertTable: InsertBasicChatConversation,
}
