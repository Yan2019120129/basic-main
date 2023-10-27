package bill

import (
	"basic/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

// TypesList 账单类型
func TypesList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make([]map[string]interface{}, 0)
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	rds := cache.RedisPool.Get()
	defer rds.Close()

	for k, v := range models.UserBillTypeNameMap {
		data = append(data, map[string]interface{}{"label": locales.Manager.GetAdminLocales(rds, adminId, "zh-CN", v), "value": k})
	}

	body.SuccessJSON(w, data)
}
