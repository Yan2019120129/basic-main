package assets

import (
	"basic/models"
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

// Options 资产Options
func Options(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)

	rds := cache.RedisPool.Get()
	defer rds.Close()

	data := make([]map[string]any, 0)
	models.NewAssets(nil).Field("id", "name", "type").
		AndWhere("admin_id=?", settingAdminId).AndWhere("status=?", models.AssetsStatusActivate).
		Query(func(rows *sql.Rows) {
			var id, assetsType int64
			var name string
			_ = rows.Scan(&id, &name, &assetsType)
			data = append(data, map[string]any{
				"label": locales.Manager.GetAdminLocales(rds, settingAdminId, "zh-CN", name), "value": id,
			})
		})
	body.SuccessJSON(w, data)
}
