package dictionary

import (
	"basic/models"
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type downloadData struct {
	Alias string `json:"alias"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	Field string `json:"field"`
	Value string `json:"value"`
}

func Download(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	data := make([]*downloadData, 0)

	// 系统管理员字典
	models.NewLangDictionary(nil).
		Field("alias", "name", "field", "value").
		AndWhere("admin_id=0").AndWhere("alias=?", "zh-CN").Query(func(rows *sql.Rows) {
		temp := new(downloadData)
		temp.Type = "0"
		_ = rows.Scan(&temp.Alias, &temp.Name, &temp.Field, &temp.Value)
		data = append(data, temp)
	})

	models.NewLangDictionary(nil).
		Field("alias", "type", "name", "field", "value").
		AndWhere("admin_id=?", adminId).AndWhere("alias=?", "zh-CN").Query(func(rows *sql.Rows) {
		temp := new(downloadData)
		_ = rows.Scan(&temp.Alias, &temp.Type, &temp.Name, &temp.Field, &temp.Value)
		data = append(data, temp)
	})

	body.SuccessJSON(w, data)
}
