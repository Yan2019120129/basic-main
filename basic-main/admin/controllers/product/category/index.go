package category

import (
	"basic/models"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type indexParams struct {
	ParentId   int64                  `json:"parent_id"`
	AdminName  string                 `json:"admin_name"`
	Name       string                 `json:"name"`
	Status     int64                  `json:"status"`
	Type       int64                  `json:"type"`
	Recommend  int64                  `json:"recommend"`
	DateTime   *define.RangeTimeParam `json:"updated_at"`
	Pagination *define.Pagination     `json:"pagination"`
}

type indexData struct {
	models.ProductCategoryAttrs
	Data       string `json:"data"`        //数据
	NameLang   string `json:"nameLang"`    //字典名称
	AdminName  string `json:"admin_name"`  //管理员名称
	ParentName string `json:"parent_name"` //父级名称
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
	model := models.NewProductCategory(nil)
	model.Db.AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("status>?", models.ProductCategoryStatusDelete)

	define.NewFilterEmpty(model.Db).
		String("name like ?", "%"+params.Name+"%").
		Int64("status=?", params.Status).
		Int64("type=?", params.Type).
		Int64("recommend=?", params.Recommend).
		Int64("parent_id=?", params.ParentId).
		RangeTime("updated_at between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)

	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	data := make([]*indexData, 0)
	rds := cache.RedisPool.Get()
	defer rds.Close()
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.ParentId, &tmp.AdminId, &tmp.Type, &tmp.Name, &tmp.Image, &tmp.Sort, &tmp.Status, &tmp.Recommend, &tmp.Data, &tmp.UpdatedAt, &tmp.CreatedAt)
		adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}
		if tmp.ParentId > 0 {
			parentInfo := models.NewProductCategory(nil).AndWhere("id=?", tmp.ParentId).FindOne()
			if parentInfo != nil {
				tmp.ParentName = locales.Manager.GetAdminLocales(rds, settingAdminId, "zh-CN", parentInfo.Name)
			}
		}

		tmp.NameLang = locales.Manager.GetAdminLocales(rds, settingAdminId, "zh-CN", tmp.Name)
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
