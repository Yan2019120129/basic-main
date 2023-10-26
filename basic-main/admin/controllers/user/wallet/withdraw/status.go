package withdraw

import (
	"basic/models"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type statusParams struct {
	Id     int64  `json:"id" validate:"required"`
	Status int64  `json:"status" validate:"required,oneof=-1 20"`
	Data   string `json:"data"`
}

func Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(statusParams)
	_ = body.ReadJSON(r, params)

	//  参数验证
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 是否存在订单
	walletOrderInfo := models.NewUserWalletOrder(nil).AndWhere("id=?", params.Id).AndWhere("status=?", models.WalletOrderStatusPending).FindOne()
	if walletOrderInfo == nil {
		body.ErrorJSON(w, "订单不存在", -1)
		return
	}

	userInfo := models.NewUser(nil).AndWhere("id=?", walletOrderInfo.UserId).FindOne()
	if userInfo == nil {
		body.ErrorJSON(w, "用户不存在", -1)
		return
	}

	//  实例化模型
	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	model := models.NewUserWalletOrder(tx)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	nowTime := time.Now()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("data=?", params.Data).
		Int64("status=?", params.Status).
		Float64("balance=?", userInfo.Money).
		Int64("updated_at=?", nowTime.Unix())

	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id=?", params.Id).AndWhere("type=?", models.WalletOrderTypeWithdraw).Update()
	if err != nil {
		panic(err)
	}

	// 如果拒绝，那么新增用户金额
	if params.Status == models.WalletOrderStatusRefuse {
		err = models.UserDeposit(tx, userInfo, walletOrderInfo.Id, models.UserBillTypeWithdrawRefuse, walletOrderInfo.Money)
		if err != nil {
			panic(err)
		}
	}

	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
