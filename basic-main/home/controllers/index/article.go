package index

import (
	"basic/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/utils/body"
)

type articleParams struct {
	Field string `json:"field" validate:"required"`
}

// Article 文章内容
func Article(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(articleParams)
	_ = body.ReadJSON(r, params)
	adminId := models.NewAdminUser(nil).GetDomainAdminId(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)

	rds := cache.RedisPool.Get()
	defer rds.Close()

	data := locales.Manager.GetAdminLocales(rds, settingAdminId, r.Header.Get("Accept-Language"), params.Field)
	body.SuccessJSON(w, data)
}
