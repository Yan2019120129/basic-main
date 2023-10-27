package product

import (
	"basic/models"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type productItem struct {
	Id     int64                  `json:"id"`
	Images string                 `json:"images"`
	Name   string                 `json:"name"`
	Money  float64                `json:"money"`
	Type   int64                  `json:"type"`
	Sales  int64                  `json:"sales"`
	Nums   int64                  `json:"nums"`
	Used   int64                  `json:"used"`
	Total  int64                  `json:"total"`
	Data   map[string]interface{} `json:"data"`
}

type indexParams struct {
	CategoryId int64 `json:"category_id"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	rds := cache.RedisPool.Get()
	defer rds.Close()
	acceptLanguage := r.Header.Get("Accept-Language")
	claims := router.TokenManager.GetContextClaims(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)

	data := make([]*productItem, 0)
	model := models.NewProduct(nil).Field("id", "images", "name", "money", "type", "sales", "nums", "used", "total", "data")

	// 判断是否携带分类ID
	if params.CategoryId > 0 {
		categoryModel := models.NewProductCategory(nil)
		categoryModel.AndWhere("id=?", params.CategoryId).AndWhere("admin_id=?", settingAdminId).AndWhere("status>?", models.ProductCategoryStatusDelete)
		categoryInfo := categoryModel.FindOne()
		if categoryInfo != nil {
			model.AndWhere("category_id=?", categoryInfo.Id)
		}
	}

	model.AndWhere("admin_id=?", settingAdminId).AndWhere("status>?", models.ProductStatusDelete).
		OrderBy("sort asc").Query(func(rows *sql.Rows) {
		temp := new(productItem)
		var oldData string
		_ = rows.Scan(&temp.Id, &temp.Images, &temp.Name, &temp.Money, &temp.Type, &temp.Sales, &temp.Nums, &temp.Used, &temp.Total, &oldData)
		_ = json.Unmarshal([]byte(oldData), &temp.Data)
		temp.Name = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, temp.Name)
		data = append(data, temp)
	})

	body.SuccessJSON(w, data)
}
