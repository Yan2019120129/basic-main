package lang

import (
	"basic/models"
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

// LanguageList 语言列表
func LanguageList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	data := make([]map[string]any, 0)
	models.NewLang(nil).Field("alias", "name").
		AndWhere("admin_id=?", adminId).AndWhere("status=?", models.LangStatusActivate).
		Query(func(rows *sql.Rows) {
			var alias string
			var name string
			_ = rows.Scan(&alias, &name)
			data = append(data, map[string]any{
				"label": name, "value": alias,
			})
		})
	body.SuccessJSON(w, data)
}
