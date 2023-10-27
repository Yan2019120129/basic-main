package product

import (
	"basic/models"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type updateParams struct {
	Id         int64               `json:"id" validate:"required"`
	AssetsId   int64               `json:"assets_id"`
	CategoryId int64               `json:"category_id" validate:"omitempty,gt=0"`
	Name       string              `json:"name"`
	Images     []map[string]string `json:"images_list"`
	Money      float64             `json:"money" validate:"omitempty,gte=0"`
	Sort       int64               `json:"sort" validate:"omitempty,gte=0"`
	Status     int64               `json:"status" validate:"omitempty,oneof=-1 10"`
	Recommend  int64               `json:"recommend" validate:"omitempty,oneof=-1 10"`
	Sales      int64               `json:"sales" validate:"omitempty,gte=0"`
	Nums       int64               `json:"nums" validate:"omitempty,gte=-1"`
	Mode       int64               `json:"mode" validate:"omitempty,oneof=1 2"`
	Used       int64               `json:"used" validate:"omitempty,gt=0"`
	Total      int64               `json:"total" validate:"omitempty,gt=0"`
	Data       string              `json:"data"`
	Describes  string              `json:"describes"`
}

func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updateParams)
	_ = body.ReadJSON(r, params)

	//  参数验证
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}
	//  实例化模型
	model := models.NewProduct(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	images := ""
	if len(params.Images) > 0 {
		imagesByte, _ := json.Marshal(params.Images)
		images = string(imagesByte)
	}

	// 判断分类ID是否存在
	if params.CategoryId > 0 {
		categoryInfo := models.NewProductCategory(nil).AndWhere("id=?", params.CategoryId).FindOne()
		if categoryInfo == nil {
			body.ErrorJSON(w, "分类ID不存在", -1)
			return
		}
	}

	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		Int64("category_id=?", params.CategoryId).
		String("name=?", params.Name).
		String("images=?", images).
		Float64("money=?", params.Money).
		Int64("sort=?", params.Sort).
		Int64("status=?", params.Status).
		Int64("recommend=?", params.Recommend).
		Int64("sales=?", params.Sales).
		Int64("nums=?", params.Nums).
		Int64("mode=?", params.Mode).
		Int64("used=?", params.Used).
		Int64("total=?", params.Total).
		String("data=?", params.Data).
		Int64("assets_id=?", params.AssetsId).
		String("describes=?", params.Describes)

	if params.AssetsId == -1 {
		model.Value("assets_id=?").Args(0)
	}

	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id = ?", params.Id).Update()
	if err != nil {
		panic(err)
	}
	body.SuccessJSON(w, "ok")
}
