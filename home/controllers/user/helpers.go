package user

import (
	"basic/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type helpersData struct {
	BackgroundImage string                   `json:"backgroundImage"`
	Helpers         []map[string]interface{} `json:"helpers"`
	Contact         []map[string]interface{} `json:"contact"`
}

// Helpers 帮助中心
func Helpers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	rds := cache.RedisPool.Get()
	defer rds.Close()

	adminSettingList := models.NewAdminSetting(nil).GetAdminFieldAllString(settingAdminId)
	helpersList := models.AdminSettingValueToMapInterfaces(adminSettingList["helpers"])
	contactList := models.AdminSettingValueToMapInterfaces(adminSettingList["contacts"])
	acceptLanguage := r.Header.Get("Accept-Language")

	// 语言问题1
	for _, helper := range helpersList {
		helper["field"] = helper["content"]
		helper["title"] = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, helper["title"].(string))
		helper["content"] = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, helper["content"].(string))
	}
	// 语言问题2
	for _, contact := range contactList {
		contact["desc"] = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, contact["desc"].(string))
	}

	body.SuccessJSON(w, &helpersData{
		BackgroundImage: adminSettingList["service_image"],
		Helpers:         helpersList,
		Contact:         contactList,
	})
}
