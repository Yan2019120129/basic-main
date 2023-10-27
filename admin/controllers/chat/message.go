package chat

import (
	"basic/models"
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type messageParams struct {
	ConversationId int64 `json:"conversationId" validate:"required"`
}

type MessageData struct {
	Id       int64  `json:"id"`
	UserId   int64  `json:"userId"`
	Avatar   string `json:"avatar"`
	NickName string `json:"nickname"`
	Type     int64  `json:"type"`
	Data     string `json:"data"`
	Time     int64  `json:"time"`
}

// Message 会话消息内容
func Message(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(messageParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	rds := cache.RedisPool.Get()
	defer rds.Close()

	claims := router.TokenManager.GetHeaderClaims(rds, r)
	if claims == nil {
		body.ErrorJSON(w, "Error Token ServiceInfo", -1)
		return
	}

	onlineInfo := models.NewUser(nil).AndWhere("admin_id=?", claims.AdminId).AndWhere("id=?", claims.UserId).AndWhere("status=?", models.UserStatusActivate).FindOne()
	if onlineInfo == nil {
		body.ErrorJSON(w, "Error OnlineInfo", -1)
		return
	}

	conversationInfo := models.NewChatConversation(nil).AndWhere("id=?", params.ConversationId).AndWhere("user_id=?", onlineInfo.Id).FindOne()
	if conversationInfo == nil {
		body.ErrorJSON(w, "管理员会话不存在", -1)
		return
	}

	data := make([]*MessageData, 0)
	models.NewChatConversationMessage(nil).Table("chat_conversation_message as msg").
		Field("msg.id", "user.id", "user.avatar", "user.nickname", "msg.type", "msg.data", "msg.created_at").
		LeftJoin("user", "user.id=msg.sender_id").AndWhere("msg.conversation_id=?", conversationInfo.ConversationId).
		OrderBy("msg.id desc").OffsetLimit(0, 100).Query(func(rows *sql.Rows) {
		attrs := new(MessageData)
		err = rows.Scan(&attrs.Id, &attrs.UserId, &attrs.Avatar, &attrs.NickName, &attrs.Type, &attrs.Data, &attrs.Time)
		if err == nil {
			data = append([]*MessageData{attrs}, data...)
		}
	})

	body.SuccessJSON(w, data)
}
