package product

import (
	"basic/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type createParams struct {
	AssetsId   int64               `json:"assets_id"`
	CategoryId int64               `json:"category_id" validate:"required"`
	Name       string              `json:"name" validate:"required"`
	Images     []map[string]string `json:"images" validate:"required"`
	Money      float64             `json:"money" validate:"gt=0"`
	Type       int64               `json:"type" validate:"required,oneof=1"`
	Data       string              `json:"data"`
	Describes  string              `json:"describes"`
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	//  验证参数
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 查询分类ID是否存在
	categoryInfo := models.NewProductCategory(nil).AndWhere("id=?", params.CategoryId).AndWhere("status=?", models.ProductCategoryStatusActivate).FindOne()
	if categoryInfo == nil {
		body.ErrorJSON(w, "分类ID不存在", -1)
		return
	}

	imagesByte, _ := json.Marshal(params.Images)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	nowTime := time.Now()
	//  模型插入数据
	_, err = models.NewProduct(nil).
		Field("admin_id", "category_id", "assets_id", "name", "images", "money", "type", "data", "describes", "updated_at", "created_at").
		Args(adminId, params.CategoryId, params.AssetsId, params.Name, string(imagesByte), params.Money, params.Type, params.Data, params.Describes, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
