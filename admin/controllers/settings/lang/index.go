package lang

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
	Name       string                 `json:"name"`
	Alias      string                 `json:"alias"`
	DateTime   *define.RangeTimeParam `json:"created_at"`
	Pagination *define.Pagination     `json:"pagination"`
}

type indexData struct {
	models.LangAttrs
	AdminName string `json:"admin_name"` //管理员名称
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取子级包括自己ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	//  实例化模型
	model := models.NewLang(nil)
	model.Db.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	define.NewFilterEmpty(model.Db).
		String("name like ?", "%"+params.Name+"%").
		String("alias like ?", "%"+params.Alias+"%").
		RangeTime("created_at between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)

	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.Name, &tmp.Alias, &tmp.Icon, &tmp.Sort, &tmp.Status, &tmp.Data, &tmp.CreatedAt)

		adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
