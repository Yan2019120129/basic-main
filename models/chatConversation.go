package models

import (
	"database/sql"
	"time"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/utils"
)

const (
	ChatConversationTypePrivateLetter = 1  //	私信
	ChatConversationStatusActivate    = 10 //	激活会话
	ChatConversationStatusDelete      = -2 //	删除会话
	ChatConversationStatusShielded    = -1 //	屏蔽会话
)

// ChatConversationAttrs 数据库模型属性
type ChatConversationAttrs struct {
	Id             int64  `json:"id"`              //主键
	ConversationId string `json:"conversation_id"` //会话ID
	UserId         int64  `json:"user_id"`         //用户ID
	ReceiverId     int64  `json:"receiver_id"`     //接收用户ID
	Type           int64  `json:"type"`            //1私聊
	Status         int64  `json:"status"`          //状态 -2删除 -1屏蔽 10正常
	Data           string `json:"data"`            //数据
	UpdatedAt      int64  `json:"updated_at"`      //更新时间
	CreatedAt      int64  `json:"created_at"`      //创建时间
}

// ChatConversation 数据库模型
type ChatConversation struct {
	define.Db
}

// NewChatConversation 创建数据库模型
func NewChatConversation(tx *sql.Tx) *ChatConversation {
	return &ChatConversation{
		database.DbPool.NewDb(tx).Table("chat_conversation"),
	}
}

// AndWhere where条件
func (c *ChatConversation) AndWhere(str string, arg ...any) *ChatConversation {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *ChatConversation) FindOne() *ChatConversationAttrs {
	attrs := new(ChatConversationAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.ConversationId, &attrs.UserId, &attrs.ReceiverId, &attrs.Type, &attrs.Status, &attrs.Data, &attrs.UpdatedAt, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *ChatConversation) FindMany() []*ChatConversationAttrs {
	data := make([]*ChatConversationAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(ChatConversationAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.ConversationId, &tmp.UserId, &tmp.ReceiverId, &tmp.Type, &tmp.Status, &tmp.Data, &tmp.UpdatedAt, &tmp.CreatedAt)
		data = append(data, tmp)
	})
	return data
}

// InitConversationFunc 初始化用户会话
func InitConversationFunc(tx *sql.Tx, senderUserId, receiverUserId int64, data string) (int64, string) {
	nowTime := time.Now()
	var conversationId = utils.NewRandom().OrderSn()
	var conversationInt64Id int64

	// 如果会话不存在， 那么创建会话
	conversationInfo := NewChatConversation(nil).
		AndWhere("user_id=?", senderUserId).AndWhere("receiver_id=?", receiverUserId).
		FindOne()
	if conversationInfo == nil {
		conversationInsertId, err := NewChatConversation(tx).
			Field("conversation_id", "user_id", "receiver_id", "type", "data", "updated_at", "created_at").
			Args(conversationId, senderUserId, receiverUserId, ChatConversationTypePrivateLetter, data, nowTime.Unix(), nowTime.Unix()).
			Insert()
		if err != nil {
			panic(err)
		}
		conversationInt64Id = conversationInsertId
	} else {
		_, err := NewChatConversation(tx).
			Value("data=?", "updated_at=?").Args(data, nowTime.Unix()).
			AndWhere("user_id=?", senderUserId).AndWhere("receiver_id=?", receiverUserId).
			Update()
		if err != nil {
			panic(err)
		}
		conversationInt64Id = conversationInfo.Id
		conversationId = conversationInfo.ConversationId
	}

	// 好友会话
	friendsConversationInfo := NewChatConversation(nil).
		AndWhere("user_id=?", receiverUserId).AndWhere("receiver_id=?", senderUserId).
		FindOne()
	if friendsConversationInfo == nil {
		_, err := NewChatConversation(tx).
			Field("conversation_id", "user_id", "receiver_id", "type", "data", "updated_at", "created_at").
			Args(conversationId, receiverUserId, senderUserId, ChatConversationTypePrivateLetter, data, nowTime.Unix(), nowTime.Unix()).
			Insert()
		if err != nil {
			panic(err)
		}
	} else {
		_, err := NewChatConversation(tx).
			Value("data=?", "updated_at=?").Args(data, nowTime.Unix()).
			AndWhere("user_id=?", receiverUserId).AndWhere("receiver_id=?", senderUserId).
			Update()
		if err != nil {
			panic(err)
		}
	}

	return conversationInt64Id, conversationId
}
