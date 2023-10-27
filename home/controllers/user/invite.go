package user

import (
	"basic/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/utils/body"
)

type inviteData struct {
	SiteName        string `json:"siteName"`
	SiteLogo        string `json:"siteLogo"`
	BackgroundImage string `json:"backgroundImage"`
}

// Invite 邀请信息
func Invite(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := models.NewAdminUser(nil).GetDomainAdminId(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	settingAdminList := models.NewAdminSetting(nil).GetAdminFieldAllString(settingAdminId)

	body.SuccessJSON(w, &inviteData{
		SiteName:        settingAdminList["site_name"],
		SiteLogo:        settingAdminList["site_logo"],
		BackgroundImage: settingAdminList["invite_image"],
	})
}
