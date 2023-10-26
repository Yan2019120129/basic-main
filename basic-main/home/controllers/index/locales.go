package index

import (
	"basic/models"
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/utils/body"
)

// Locales 语言列表
func Locales(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := models.NewAdminUser(nil).GetDomainAdminId(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)

	locales := make([]*Locale, 0)
	models.NewLangDictionary(nil).Field("field", "value").
		AndWhere("alias=?", r.Header.Get("Accept-Language")).
		AndWhere("admin_id=?", settingAdminId).AndWhere("type=?", models.LangDictionaryTypeHomeTranslate).
		Query(func(rows *sql.Rows) {
			localeTmp := new(Locale)
			_ = rows.Scan(&localeTmp.Label, &localeTmp.Value)
			locales = append(locales, localeTmp)
		})

	body.SuccessJSON(w, locales)
}
