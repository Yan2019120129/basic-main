package assets

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
	Name string `json:"name" validate:"required,max=50"`
	Icon string `json:"icon" validate:"required,max=193"`
	Type int64  `json:"type" validate:"required,oneof=1 2 3"`
	Data string `json:"data"`
}

// Create 新增资产
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)
	//  验证参数
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	//  模型插入数据
	nowTime := time.Now()
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	//	不能重复的别名
	oldCountryInfo := models.NewAssets(nil).AndWhere("admin_id=?", adminId).AndWhere("name=?", params.Name).AndWhere("status>=?", models.AssetsStatusDisabled).FindOne()
	if oldCountryInfo != nil {
		body.ErrorJSON(w, "当前资产已存在", -1)
		return
	}

	_, err = models.NewAssets(nil).
		Field("admin_id", "name", "icon", "type", "data", "created_at", "updated_at").
		Args(adminId, params.Name, params.Icon, params.Type, params.Data, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}
	body.SuccessJSON(w, "ok")
}
