package sql

import (
	"basic/tools/utils"
)

const BasicChatConversationMessageTableName = "chat_conversation_message"
const BasicChatConversationMessageTableComment = "聊天会话消息"
const CreateBasicChatConversationMessage = `CREATE TABLE ` + BasicChatConversationMessageTableName + ` (
	id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	conversation_id CHAR(32) NOT NULL COMMENT '会话ID',
	sender_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '发送者ID',
	receiver_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '接收者用户ID',
	unread TINYINT NOT NULL DEFAULT 1 COMMENT '未读 1未读 2已读',
	type TINYINT NOT NULL DEFAULT 1 COMMENT '消息类型 0起点 1文本 2图片 3语音 4视频 10富文本',
	data TEXT COMMENT '消息内容',
	extra VARCHAR(2048) NOT NULL DEFAULT '' COMMENT '扩展数据',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
	KEY message_conversation_id (conversation_id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + BasicChatConversationMessageTableComment + `';`

const InsertBasicChatConversationMessage = ``

var BasicChatConversationMessage = &utils.InitTable{
	Name:        BasicChatConversationMessageTableName,
	Comment:     BasicChatConversationMessageTableComment,
	CreateTable: CreateBasicChatConversationMessage,
	InsertTable: InsertBasicChatConversationMessage,
}
