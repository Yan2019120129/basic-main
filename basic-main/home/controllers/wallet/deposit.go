package wallet

import (
	"basic/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type depositParams struct {
	Id     int64   `json:"id" validate:"required"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Proof  string  `json:"proof"`
}

type depositData struct {
	Url string `json:"url"` //	跳转链接
}

// Deposit 充值
func Deposit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(depositParams)
	_ = body.ReadJSON(r, params)

	rds := cache.RedisPool.Get()
	defer rds.Close()

	claims := router.TokenManager.GetContextClaims(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)
	adminSettingList := models.NewAdminSetting(nil).GetAdminFieldAllString(settingAdminId)
	acceptLanguage := r.Header.Get("Accept-Language")
	userModel := models.NewUser(nil)
	userModel.AndWhere("id=?", claims.UserId)
	userInfo := userModel.FindOne()

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 判断是否充值范围
	depositRangeList := models.AdminSettingValueToMapInterface(adminSettingList["finance_deposit_range"])
	depositMinAmount := depositRangeList["min"].(float64)
	depositMaxAmount := depositRangeList["max"].(float64)
	if params.Amount < depositMinAmount || params.Amount > depositMaxAmount {
		errMsg := locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "depositRangeError")
		errMsg = strings.ReplaceAll(errMsg, "{min}", strconv.FormatFloat(depositMinAmount, 'f', 0, 64))
		errMsg = strings.ReplaceAll(errMsg, "{max}", strconv.FormatFloat(depositMaxAmount, 'f', 10, 64))
		body.ErrorJSON(w, errMsg, -1)
		return
	}

	// 如果当前用户有未完成订单，那么提示请联系客服人员
	unfulfilledOrderModel := models.NewUserWalletOrder(nil)
	unfulfilledOrderModel.AndWhere("user_id=?", claims.UserId).AndWhere("type=?", models.WalletOrderTypeDeposit).AndWhere("status=?", models.WalletOrderStatusPending)
	unfulfilledOrderInfo := unfulfilledOrderModel.FindOne()
	if unfulfilledOrderInfo != nil {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "thereAreStillOutstandingOrders"), -1)
		return
	}

	// 判断支付方式是否存在
	paymentModel := models.NewWalletPayment(nil)
	paymentModel.AndWhere("id=?", params.Id).AndWhere("admin_id=?", settingAdminId).AndWhere("status=?", models.WalletPaymentStatusActivate).AndWhere("mode=?", models.WalletPaymentModeDeposit)
	paymentInfo := paymentModel.FindOne()
	if paymentInfo == nil {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "noPermissionToExecute"), -1)
		return
	}

	// 如果是银行卡或加密货币那么需要凭证
	if paymentInfo.Type == models.WalletPaymentTypeBank || paymentInfo.Type == models.WalletPaymentTypeCryptocurrency {
		if params.Proof == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "credentialsCannotBeEmpty"), -1)
			return
		}
	}

	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	nowTime := time.Now()
	orderSn := utils.NewRandom().OrderSn()
	_, err = models.NewUserWalletOrder(tx).
		Field("order_sn", "admin_id", "user_id", "type", "payment_id", "money", "balance", "proof", "updated_at", "created_at").
		Args(orderSn, userInfo.AdminId, userInfo.Id, models.WalletOrderTypeDeposit, paymentInfo.Id, params.Amount, userInfo.Money, params.Proof, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	// 充值提示音
	// TODO...

	_ = tx.Commit()
	body.SuccessJSON(w, &depositData{
		Url: "",
	})
}
