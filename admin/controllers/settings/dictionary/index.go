package dictionary

import (
	"basic/models"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type indexParams struct {
	AdminName  string                 `json:"admin_name"`
	LangAlias  string                 `json:"alias"`
	Type       int64                  `json:"type"`
	Name       string                 `json:"name"`
	Field      string                 `json:"field"`
	Value      string                 `json:"value"`
	Data       string                 `json:"data"`
	DateTime   *define.RangeTimeParam `json:"created_at"`
	Pagination *define.Pagination     `json:"pagination"`
}

type indexData struct {
	models.LangDictionaryAttrs
	AdminName string `json:"admin_name"`
	LangName  string `json:"lang_name"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  实例化模型
	model := models.NewLangDictionary(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	model.Db.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	define.NewFilterEmpty(model.Db).
		Int64("type=?", params.Type).
		String("name like ?", "%"+params.Name+"%").
		String("field like ?", "%"+params.Field+"%").
		String("value like ?", "%"+params.Value+"%").
		String("data like ?", "%"+params.Data+"%").
		String("alias=?", params.LangAlias).
		RangeTime("created_at between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)

	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.Type, &tmp.Alias, &tmp.Name, &tmp.Field, &tmp.Value, &tmp.Data, &tmp.CreatedAt)

		// 当前管理员信息
		adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}

		// 当前语言信息
		langInfo := models.NewLang(nil).AndWhere("alias=?", tmp.Alias).FindOne()
		if langInfo != nil {
			tmp.LangName = langInfo.Name
		}
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
