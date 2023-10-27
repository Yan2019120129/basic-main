package product

import (
	"basic/models"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type detailsData struct {
	Id        int64                  `json:"id"`
	Name      string                 `json:"name"`
	Images    []map[string]string    `json:"images"`
	Describes string                 `json:"describes"`
	Data      map[string]interface{} `json:"data"`
	Money     float64                `json:"money"`
	Recommend int64                  `json:"recommend"`
	Sales     int64                  `json:"sales"`
}

type detailsParams struct {
	Id int64 `json:"id" validate:"required"`
}

func Details(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(detailsParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	acceptLanguage := r.Header.Get("Accept-Language")
	rds := cache.RedisPool.Get()
	defer rds.Close()

	claims := router.TokenManager.GetContextClaims(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)

	productModel := models.NewProduct(nil)
	productModel.AndWhere("id=?", params.Id).AndWhere("admin_id=?", settingAdminId).AndWhere("status>?", models.ProductStatusDelete)
	productInfo := productModel.FindOne()
	if productInfo == nil {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "noPermissionToExecute"), -1)
		return
	}

	images := make([]map[string]string, 0)
	if productInfo.Images != "" {
		_ = json.Unmarshal([]byte(productInfo.Images), &images)
	}

	var productData map[string]interface{}
	if productInfo.Data != "" {
		_ = json.Unmarshal([]byte(productInfo.Data), &productData)
	}

	productName := locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, productInfo.Name)
	describes := locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, productInfo.Describes)
	body.SuccessJSON(w, &detailsData{
		Id:        productInfo.Id,
		Name:      productName,
		Images:    images,
		Describes: describes,
		Data:      productData,
		Money:     productInfo.Money,
		Recommend: productInfo.Recommend,
		Sales:     productInfo.Sales,
	})
}
