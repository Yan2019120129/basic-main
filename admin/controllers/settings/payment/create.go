package payment

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
	Name        string `json:"name" validate:"required"`
	Icon        string `json:"icon"`
	Mode        int64  `json:"mode" validate:"required,oneof=1 10"`
	Type        int64  `json:"type" validate:"required,oneof=1 10 20"`
	AccountName string `json:"account_name" validate:"required"`
	AccountCode string `json:"account_code" validate:"required"`
	Sort        int64  `json:"sort"`
	Data        string `json:"data"`
	Expand      string `json:"expand"` //扩展数据
	Description string `json:"description"`
}

// Create 新增系统钱包
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	//  验证参数
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	adminId := router.TokenManager.GetContextClaims(r).AdminId
	nowTime := time.Now()

	//  模型插入数据
	_, err = models.NewWalletPayment(nil).
		Field("admin_id", "name", "icon", "mode", "type", "account_name", "account_code", "sort", "description", "data", "expand", "created_at").
		Args(adminId, params.Name, params.Icon, params.Mode, params.Type, params.AccountName, params.AccountCode, params.Sort, params.Description, params.Data, params.Expand, nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
