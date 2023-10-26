package country

import (
	"basic/models"
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

// Options 选择框参数
func Options(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	model := models.NewCountry(nil).AndWhere("admin_id=?", adminId).AndWhere("status=?", models.CountryStatusActivate)

	data := make([]map[string]any, 0)
	model.Field("id", "name").Query(func(rows *sql.Rows) {
		var name string
		var id int64
		_ = rows.Scan(&id, &name)
		data = append(data, map[string]any{"label": name, "value": id})
	})

	body.SuccessJSON(w, data)
}
