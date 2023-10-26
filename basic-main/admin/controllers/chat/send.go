package chat

import (
	"basic/models"
	"basic/module/socket"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type sendparams struct {
	ConversationId int64  `json:"conversationId"`
	Type           int64  `json:"type" validate:"required,oneof=1 2"`
	Data           string `json:"data" validate:"required"`
}

type sendData struct {
	MessageId      int64 `json:"messageId"`
	ConversationId int64 `json:"conversationId"`
}

// Send 发送消息
func Send(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(sendparams)
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

	//	token用户信息
	tokenUserInfo := models.NewUser(nil).AndWhere("admin_id=?", claims.AdminId).AndWhere("id=?", claims.UserId).AndWhere("status=?", models.UserStatusActivate).FindOne()
	if tokenUserInfo == nil {
		body.ErrorJSON(w, "Error OnlineInfo", -1)
		return
	}

	// 是否有会话ID
	var senderUserId, receiverUserId int64
	conversationInfo := models.NewChatConversation(nil).AndWhere("id=?", params.ConversationId).FindOne()
	if conversationInfo == nil {
		adminOnlineInfo := models.NewUser(nil).AndWhere("admin_id=?", claims.AdminId).AndWhere("type=?", models.UserTypeOnline).AndWhere("status=?", models.UserStatusActivate).FindOne()
		if adminOnlineInfo == nil || adminOnlineInfo.Id == tokenUserInfo.Id {
			body.ErrorJSON(w, "Error AdminOnlineInfo", -1)
			return
		}
		senderUserId = tokenUserInfo.Id
		receiverUserId = adminOnlineInfo.Id
	} else {
		senderUserId = conversationInfo.UserId
		receiverUserId = conversationInfo.ReceiverId
	}

	nowTime := time.Now()
	tx := database.DbPool.GetTx()
	defer tx.Rollback()
	conversationData := params.Data
	if params.Type == models.ChatConversationTypeImage {
		conversationData = "[图片]"
	}

	// 初始化会话
	conversationInt64Id, conversationId := models.InitConversationFunc(tx, senderUserId, receiverUserId, conversationData)
	messageId, err := models.NewChatConversationMessage(tx).
		Field("conversation_id", "sender_id", "receiver_id", "type", "data", "updated_at", "created_at").
		Args(conversationId, senderUserId, receiverUserId, params.Type, params.Data, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	// 发送websocket消息
	msg := &MessageData{Id: messageId, UserId: tokenUserInfo.Id, Avatar: tokenUserInfo.Avatar, NickName: tokenUserInfo.Nickname, Type: params.Type, Data: params.Data, Time: nowTime.Unix()}
	socket.SocketInstance.SendChatMessage(receiverUserId, msg)

	_ = tx.Commit()
	body.SuccessJSON(w, &sendData{MessageId: messageId, ConversationId: conversationInt64Id})
}
