package country

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
	LangName   string                 `json:"lang_name"`
	Name       string                 `json:"name"`
	Alias      string                 `json:"alias"`
	Iso1       string                 `json:"iso1"` //	ISO3166-1
	Code       string                 `json:"code"`
	DateTime   *define.RangeTimeParam `json:"created_at"`
	Pagination *define.Pagination     `json:"pagination"`
}

type indexData struct {
	models.CountryAttrs
	AdminName string `json:"admin_name"` //管理员名称
	LangName  string `json:"lang_name"`  //语言名称
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  实例化模型
	model := models.NewCountry(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	model.Db.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	define.NewFilterEmpty(model.Db).
		String("name like ?", "%"+params.Name+"%").
		String("alias like ?", "%"+params.Alias+"%").
		String("iso1 like ?", "%"+params.Iso1+"%").
		String("code like ?", "%"+params.Code+"%").
		RangeTime("created_at between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)

	if params.LangName != "" {
		langNameIds := models.NewLang(nil).Field("id").AndWhere("name like ?", "%"+params.LangName+"%").ColumnString()
		if len(langNameIds) == 0 {
			langNameIds = append(langNameIds, "-1")
		}
		model.Db.AndWhere("lang_id in (" + strings.Join(langNameIds, ",") + ")")
	}
	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.LangId, &tmp.Name, &tmp.Alias, &tmp.Iso1, &tmp.Icon, &tmp.Code, &tmp.Sort, &tmp.Status, &tmp.Data, &tmp.CreatedAt)
		// 当前管理员信息
		adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}

		// 当前语言信息
		langInfo := models.NewLang(nil).AndWhere("id=?", tmp.LangId).FindOne()
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
