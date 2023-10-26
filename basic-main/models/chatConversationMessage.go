package models

import (
	"database/sql"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

const (
	ChatConversationMessageUnread = 1 //	消息未读
	ChatConversationMessageRead   = 2 //	消息已读
	ChatConversationTypeText      = 1 //	文本
	ChatConversationTypeImage     = 2 //	图片
)

// ChatConversationMessageAttrs 数据库模型属性
type ChatConversationMessageAttrs struct {
	Id             int64  `json:"id"`              //主键
	ConversationId string `json:"conversation_id"` //会话ID
	SenderId       int64  `json:"sender_id"`       //发送者ID
	ReceiverId     int64  `json:"receiver_id"`     //接收者ID
	Unread         int64  `json:"unread"`          //未读 1未读 2已读
	Type           int64  `json:"type"`            //消息类型 1文本 2图片 3语音 4视频 10富文本
	Data           string `json:"data"`            //消息内容
	Extra          string `json:"extra"`           //扩展数据
	UpdatedAt      int64  `json:"updated_at"`      //更新时间
	CreatedAt      int64  `json:"created_at"`      //创建时间
}

// ChatConversationMessage 数据库模型
type ChatConversationMessage struct {
	define.Db
}

// NewChatConversationMessage 创建数据库模型
func NewChatConversationMessage(tx *sql.Tx) *ChatConversationMessage {
	return &ChatConversationMessage{
		database.DbPool.NewDb(tx).Table("chat_conversation_message"),
	}
}

// AndWhere where条件
func (c *ChatConversationMessage) AndWhere(str string, arg ...any) *ChatConversationMessage {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *ChatConversationMessage) FindOne() *ChatConversationMessageAttrs {
	attrs := new(ChatConversationMessageAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.ConversationId, &attrs.SenderId, &attrs.ReceiverId, &attrs.Unread, &attrs.Type, &attrs.Data, &attrs.Extra, &attrs.UpdatedAt, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *ChatConversationMessage) FindMany() []*ChatConversationMessageAttrs {
	data := make([]*ChatConversationMessageAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(ChatConversationMessageAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.ConversationId, &tmp.SenderId, &tmp.ReceiverId, &tmp.Unread, &tmp.Type, &tmp.Data, &tmp.Extra, &tmp.UpdatedAt, &tmp.CreatedAt)
		data = append(data, tmp)
	})
	return data
}
