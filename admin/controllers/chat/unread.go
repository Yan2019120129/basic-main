package chat

import (
	"basic/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type unreadParasm struct {
	ConversationId int64 `json:"conversationId" validate:"required"`
}

// Unread 未读消息
func Unread(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(unreadParasm)
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

	_, err = models.NewChatConversationMessage(nil).Value("unread=?").Args(models.ChatConversationMessageRead).
		AndWhere("receiver_id=?", conversationInfo.UserId).AndWhere("sender_id=?", conversationInfo.ReceiverId).Update()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
