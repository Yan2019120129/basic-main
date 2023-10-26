package level

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
	UpdatedAt  *define.RangeTimeParam `json:"updated_at"`
	Pagination *define.Pagination     `json:"pagination"` //	分页
}

type indexItems struct {
	AdminId   int64  `json:"admin_id"`   //管理员ID
	AdminName string `json:"admin_name"` //管理员名称
	models.UserLevelAttrs
}

// Index 等级列表
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取请求管理员的adminId 获取子级包括自己ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	model := models.NewUserLevel(nil)
	model.AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("status<>?", models.UserLevelStatusDelete)
	define.NewFilterEmpty(model.Db).
		String("name like ?", "%"+params.Name+"%").
		RangeTime("created_at between ? and ?", params.UpdatedAt, location).
		Pagination(params.Pagination)

	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	data := make([]*indexItems, 0)
	model.Field("id", "admin_id", "name", "icon", "level", "money", "days", "status", "data", "updated_at", "created_at").
		Query(func(rows *sql.Rows) {
			tmp := new(indexItems)
			_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.Name, &tmp.Icon, &tmp.Level, &tmp.Money, &tmp.Days, &tmp.Status, &tmp.Data, &tmp.UpdatedAt, &tmp.CreatedAt)
			adminModel := models.NewAdminUser(nil)
			adminModel.AndWhere("id=?", tmp.AdminId)
			adminInfoTmp := adminModel.FindOne()
			if adminInfoTmp != nil {
				tmp.AdminName = adminInfoTmp.UserName
			}
			data = append(data, tmp)
		})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
