package deposit

import (
	"basic/models"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type deleteParams struct {
	Ids []int64 `json:"id" validate:"required"`
}

func Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(deleteParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	nowTime := time.Now()
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	for _, v := range params.Ids {
		model := models.NewUserWalletOrder(nil)
		if adminId != models.AdminUserSupermanId {
			adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
			model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
		}
		_, err = model.Value("status=?", "updated_at=?").Args(models.WalletOrderStatusDelete, nowTime.Unix()).
			AndWhere("id=?", v).AndWhere("type=?", models.WalletOrderTypeDeposit).Update()
		if err != nil {
			panic(err)
		}
	}

	body.SuccessJSON(w, "ok")
}
