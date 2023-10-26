package chat

import (
	"basic/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type infoData struct {
	Id        int64                  `json:"id"`
	UserName  string                 `json:"username"`
	NickName  string                 `json:"nickname"`
	Avatar    string                 `json:"avatar"`
	Type      int64                  `json:"type"`
	Translate map[string]interface{} `json:"translate"`
}

// Info 在线客服信息
func Info(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rds := cache.RedisPool.Get()
	defer rds.Close()

	currentLang := r.URL.Query().Get("lang")
	if currentLang == "" {
		currentLang = "zh-CN"
	}

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

	data := &infoData{
		Id: onlineInfo.Id, UserName: onlineInfo.UserName, NickName: onlineInfo.Nickname,
		Avatar: onlineInfo.Avatar, Type: onlineInfo.Type,
	}

	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(onlineInfo.AdminId)
	data.Translate = map[string]interface{}{
		"enterText": locales.Manager.GetAdminLocales(rds, settingAdminId, currentLang, "pleaseEnterContent"),
	}

	body.SuccessJSON(w, data)
}
