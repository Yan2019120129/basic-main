package deposit

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type createParams struct {
	UserName string  `json:"username"`
	Money    float64 `json:"money"`
	Proof    string  `json:"proof"`
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

	// 查询用户是否存在
	userInfo := models.NewUser(nil).AndWhere("username=?", params.UserName).FindOne()
	if userInfo == nil {
		body.ErrorJSON(w, "用户不存在", -1)
		return
	}

	// 是否有支付方式
	paymentInfo := models.NewWalletPayment(nil).AndWhere("mode=?", models.WalletPaymentModeDeposit).FindOne()
	if paymentInfo == nil {
		body.ErrorJSON(w, "没有充值方式", -1)
		return
	}

	//  获取管理员ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	if adminId != models.AdminUserSupermanId && adminId != userInfo.AdminId {
		body.ErrorJSON(w, "权限不足", -1)
		return
	}

	nowTime := time.Now()
	orderSn := utils.NewRandom().OrderSn()
	_, err = models.NewUserWalletOrder(nil).
		Field("order_sn", "admin_id", "user_id", "user_type", "type", "payment_id", "money", "proof", "updated_at", "created_at").
		Args(orderSn, userInfo.AdminId, userInfo.Id, userInfo.Type, models.WalletOrderTypeDeposit, paymentInfo.Id, params.Money, params.Proof, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
