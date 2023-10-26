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
)

type ConversationTmp struct {
	Id       int64  `json:"id"`
	UserId   int64  `json:"userId"`
	Avatar   string `json:"avatar"`
	UserName string `json:"username"`
	NickName string `json:"nickname"`
	Data     string `json:"data"`
	Online   bool   `json:"online"`
	Ip4      string `json:"ip4"`
	Address  string `json:"address"`
	Unread   int64  `json:"unread"`
	Time     int64  `json:"time"`
}

// Conversation 会话列表
func Conversation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	data := make([]*ConversationTmp, 0)
	models.NewChatConversation(nil).Table("chat_conversation cvn").
		LeftJoin("user", "user.id=cvn.receiver_id").
		Field("cvn.id", "user.id", "user.avatar", "user.username", "user.nickname", "INET_NTOA(ip4)", "cvn.data", "cvn.updated_at").
		OrderBy("cvn.updated_at desc").OffsetLimit(0, 50).
		AndWhere("cvn.user_id=?", onlineInfo.Id).AndWhere("cvn.status=?", models.UserStatusActivate).Query(func(rows *sql.Rows) {
		attrs := new(ConversationTmp)
		err := rows.Scan(&attrs.Id, &attrs.UserId, &attrs.Avatar, &attrs.UserName, &attrs.NickName, &attrs.Ip4, &attrs.Data, &attrs.Time)
		if err == nil {
			attrs.Online = socket.SocketInstance.IsOnline(attrs.UserId)
			ip2location, _ := models.GetIp2Location(attrs.Ip4)
			if ip2location != nil {
				attrs.Address = ip2location.Country_long + "." + ip2location.Region + "." + ip2location.City
			}

			attrs.Unread = models.NewChatConversationMessage(nil).AndWhere("receiver_id=?", onlineInfo.Id).AndWhere("sender_id=?", attrs.UserId).AndWhere("unread=?", models.ChatConversationMessageUnread).Count()
			data = append(data, attrs)
		}
	})

	if len(data) == 0 {
		adminOnlineInfo := models.NewUser(nil).AndWhere("admin_id=?", claims.AdminId).AndWhere("type=?", models.UserTypeOnline).AndWhere("status=?", models.UserStatusActivate).FindOne()
		if adminOnlineInfo == nil {
			body.ErrorJSON(w, "Error AdminOnlineInfo", -1)
			return
		}

		//	判断当前用户跟管理用户是否相等
		if adminOnlineInfo.Id != claims.UserId {
			data = append(data, &ConversationTmp{Id: 0, UserId: adminOnlineInfo.Id, Avatar: adminOnlineInfo.Avatar, UserName: adminOnlineInfo.UserName, NickName: adminOnlineInfo.Nickname, Data: adminOnlineInfo.Data})
		}
	}

	body.SuccessJSON(w, data)
}
