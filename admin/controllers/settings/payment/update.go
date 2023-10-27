package payment

import (
	"basic/models"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type updateParams struct {
	Id          int64  `json:"id" validate:"required"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Mode        int64  `json:"mode" validate:"omitempty,oneof=1 10"`
	Type        int64  `json:"type" validate:"omitempty,oneof=1 10 20"`
	AccountName string `json:"account_name"`
	AccountCode string `json:"account_code"`
	Status      int64  `json:"status"`
	Sort        int64  `json:"sort"`
	Description string `json:"description"`
	Data        string `json:"data"`
	Expand      string `json:"expand"` //扩展数据
}

// Update 系统钱包更新
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
	nowTime := time.Now()
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	model := models.NewWalletPayment(nil)
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("name=?", params.Name).
		String("icon=?", params.Icon).
		Int64("mode=?", params.Mode).
		Int64("type=?", params.Type).
		String("account_name=?", params.AccountName).
		String("account_code=?", params.AccountCode).
		Int64("sort=?", params.Sort).
		Int64("status=?", params.Status).
		String("description=?", params.Description).
		String("data=?", params.Data).
		String("expand=?", params.Expand).
		Int64("updated_at=?", nowTime.Unix())

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
