package product

import (
	"basic/models"
	"database/sql"
	"encoding/json"
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
	CategoryId int64                  `json:"category_id"`
	AdminName  string                 `json:"admin_name"`
	AssetsId   int64                  `json:"assets_id"`
	Name       string                 `json:"name"`
	Status     int64                  `json:"status"`
	Recommend  int64                  `json:"recommend"`
	DateTime   *define.RangeTimeParam `json:"updated_at"`
	Pagination *define.Pagination     `json:"pagination"`
}

type indexData struct {
	Data string `json:"data"` //	数据
	models.ProductAttrs
	AdminName    string              `json:"admin_name"`    //管理员名称
	CategoryName string              `json:"category_name"` //分类ID
	CategoryType int64               `json:"category_type"` //分类类型
	AssetsName   string              `json:"assets_name"`   //	资产名称
	NameLang     string              `json:"nameLang"`      //	产品名称
	ImagesList   []map[string]string `json:"images_list"`   // 产品图片
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
	model := models.NewProduct(nil)
	model.Db.AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("status>?", models.ProductStatusDelete)

	define.NewFilterEmpty(model.Db).
		String("name like ?", "%"+params.Name+"%").
		Int64("status=?", params.Status).
		Int64("recommend=?", params.Recommend).
		Int64("category_id=?", params.CategoryId).
		Int64("assets_id=?", params.AssetsId).
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
		var productDataStr string
		imagesList := make([]map[string]string, 0)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.CategoryId, &tmp.AssetsId, &tmp.Name, &tmp.Images, &tmp.Money, &tmp.Type, &tmp.Sort, &tmp.Status, &tmp.Recommend, &tmp.Sales, &tmp.Nums, &tmp.Used, &tmp.Total, &productDataStr, &tmp.Describes, &tmp.UpdatedAt, &tmp.CreatedAt)
		adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}
		if tmp.CategoryId > 0 {
			categoryInfo := models.NewProductCategory(nil).AndWhere("id=?", tmp.CategoryId).FindOne()
			if categoryInfo != nil {
				tmp.CategoryName = locales.Manager.GetAdminLocales(rds, settingAdminId, "zh-CN", categoryInfo.Name)
				tmp.CategoryType = categoryInfo.Type
			}
		}

		tmp.AssetsName = "默认资产"
		assetsInfo := models.NewAssets(nil).AndWhere("id=?", tmp.AssetsId).FindOne()
		if assetsInfo != nil {
			tmp.AssetsName = locales.Manager.GetAdminLocales(rds, settingAdminId, "zh-CN", assetsInfo.Name)
		}

		if tmp.Images != "" {
			_ = json.Unmarshal([]byte(tmp.Images), &imagesList)
		}
		tmp.ImagesList = imagesList

		tmp.NameLang = locales.Manager.GetAdminLocales(rds, settingAdminId, "zh-CN", tmp.Name)
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
