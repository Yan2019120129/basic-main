package category

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type createParams struct {
	ParentId int64  `json:"parent_id" validate:"omitempty,gt=0"`
	Image    string `json:"image" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Type     int64  `json:"type" validate:"required,oneof=1 2 10"`
	Data     string `json:"data"`
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

	// 查询父级ID是否存在
	if params.ParentId > 0 {
		parentInfo := models.NewProductCategory(nil).AndWhere("id=?", params.ParentId).AndWhere("status=?", models.ProductCategoryStatusActivate).FindOne()
		if parentInfo == nil {
			body.ErrorJSON(w, "父级ID不存在", -1)
			return
		}
	}

	adminId := router.TokenManager.GetContextClaims(r).AdminId
	nowTime := time.Now()
	//  模型插入数据
	_, err = models.NewProductCategory(nil).
		Field("parent_id", "admin_id", "type", "name", "image", "data", "updated_at", "created_at").
		Args(params.ParentId, adminId, params.Type, params.Name, params.Image, params.Data, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
