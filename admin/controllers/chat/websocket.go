package chat

import (
	"basic/models"
	"basic/module/socket"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
)

const (
	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Resolve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Websocket 用户连接
func Websocket(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	connKey := utils.NewRandom().OrderSn()
	conn.SetReadLimit(512)
	claims := GetHeaderClaims(r)
	if claims == nil {
		return
	}

	//	所属管理员的客服ID
	adminOnlineInfo := models.NewUser(nil).AndWhere("admin_id=?", claims.AdminId).AndWhere("type=?", models.UserTypeOnline).AndWhere("status=?", models.UserStatusActivate).FindOne()
	if adminOnlineInfo == nil {
		return
	}

	// 上线通知
	socket.SocketInstance.SendChatOnlineStatus(adminOnlineInfo.Id, claims.UserId, true)

	socket.SocketInstance.SetWebsocketConn(claims.UserId, connKey, &socket.WebsocketControl{Conn: conn})
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// 下线通知
			socket.SocketInstance.SendChatOnlineStatus(adminOnlineInfo.Id, claims.UserId, false)
			break
		}

		// 接收心跳
		if string(msg) == "ping" {
			_ = conn.SetReadDeadline(time.Now().Add(pongWait))
			continue
		}
	}
}

// GetHeaderClaims 为了释放close
func GetHeaderClaims(r *http.Request) *router.Claims {
	rds := cache.RedisPool.Get()
	defer rds.Close()

	return router.TokenManager.GetHeaderClaims(rds, r)
}
