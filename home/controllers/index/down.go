package index

import (
	"basic/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/utils/body"
)

// Download 下载文件
func Download(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := models.NewAdminUser(nil).GetDomainAdminId(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)

	downFile := models.AdminSettingValueToMapInterface(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_down"))
	body.SuccessJSON(w, downFile)
}
