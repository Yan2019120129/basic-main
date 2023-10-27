package chat

import (
	"basic/models"
	"basic/module/socket"
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type conversationInfoParams struct {
	UserId int64 `json:"userId" validate:"required"`
}

// ConversationInfo 获取会话信息
func ConversationInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(conversationInfoParams)
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

	onlineInfo := models.NewUser(nil).AndWhere("id=?", claims.UserId).AndWhere("status=?", models.UserStatusActivate).FindOne()
	if onlineInfo == nil {
		body.ErrorJSON(w, "Error OnlineInfo", -1)
		return
	}

	attrs := new(ConversationTmp)
	models.NewChatConversation(nil).Table("chat_conversation cvn").
		LeftJoin("user", "user.id=cvn.receiver_id").
		Field("cvn.id", "user.id", "user.avatar", "user.username", "user.nickname", "INET_NTOA(ip4)", "cvn.data", "cvn.updated_at").
		AndWhere("cvn.user_id=?", onlineInfo.Id).AndWhere("cvn.receiver_id=?", params.UserId).AndWhere("cvn.status=?", models.UserStatusActivate).QueryRow(func(row *sql.Row) {
		err = row.Scan(&attrs.Id, &attrs.UserId, &attrs.Avatar, &attrs.UserName, &attrs.NickName, &attrs.Ip4, &attrs.Data, &attrs.Time)
		if err == nil {
			attrs.Online = socket.SocketInstance.IsOnline(attrs.UserId)
			ip2location, _ := models.GetIp2Location(attrs.Ip4)
			if ip2location != nil {
				attrs.Address = ip2location.Country_long + "." + ip2location.Region + "." + ip2location.City
			}

			attrs.Unread = models.NewChatConversationMessage(nil).AndWhere("receiver_id=?", onlineInfo.Id).AndWhere("sender_id=?", attrs.UserId).AndWhere("unread=?", models.ChatConversationMessageUnread).Count()
		}
	})

	body.SuccessJSON(w, attrs)
}
