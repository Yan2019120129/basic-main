package category

import (
	"basic/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

// Types 所有类型
func Types(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	rds := cache.RedisPool.Get()
	defer rds.Close()

	data := models.ProductCategoryTypes(rds, int64(0), settingAdminId, "")
	body.SuccessJSON(w, data)
}
