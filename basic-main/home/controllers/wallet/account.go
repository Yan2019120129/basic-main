package wallet

import (
	"basic/models"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type accountParams struct {
	Id         int64  `json:"id" validate:"required"` //	绑定类型ID
	Name       string `json:"name"`                   //	建设银行｜波场链
	RealName   string `json:"real_name"`              //	真实姓名｜USDT
	CardNumber string `json:"card_number"`            //	卡号｜地址
	Address    string `json:"address"`                //	银行地址
}

// Account 钱包账户
func Account(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(accountParams)
	_ = body.ReadJSON(r, params)
	claims := router.TokenManager.GetContextClaims(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	rds := cache.RedisPool.Get()
	defer rds.Close()
	acceptLanguage := r.Header.Get("Accept-Language")

	// 判断支付方式是否存在
	paymentModel := models.NewWalletPayment(nil)
	paymentModel.AndWhere("id=?", params.Id).AndWhere("admin_id=?", settingAdminId).AndWhere("status=?", models.WalletPaymentStatusActivate).AndWhere("mode=?", models.WalletPaymentModeWithdraw)
	paymentInfo := paymentModel.FindOne()
	if paymentInfo == nil {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "paymentMethodDoesNotExist"), -1)
		return
	}

	// 分开验证
	switch paymentInfo.Type {
	case models.WalletPaymentTypeBank:
		if params.Name == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "bankNameCannotBeEmpty"), -1)
			return
		}

		if params.RealName == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "bankRealNameCannotBeEmpty"), -1)
			return
		}

		if params.CardNumber == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "bankCardNumberCannotBeEmpty"), -1)
			return
		}

		if params.Address == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "bankAddressCannotBeEmpty"), -1)
			return
		}
	case models.WalletPaymentTypeCryptocurrency:
		if params.Name == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "chainNameCannotBeEmpty"), -1)
			return
		}

		if params.RealName == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "contractNameCannotBeEmpty"), -1)
			return
		}

		if params.CardNumber == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "contractAddressCannotBeEmpty"), -1)
			return
		}
	default:
		panic("home/controllers/wallet/account.go | 没有该类型账户")
	}

	//	单类型数量不超过多少
	accountNums, _ := strconv.ParseInt(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "financial_wallet_num"), 10, 64)
	userPaymentNums := models.NewUserWalletAccount(nil).AndWhere("user_id=?", claims.UserId).AndWhere("payment_id=?", paymentInfo.Id).AndWhere("status>?", models.UserWalletAccountStatusDelete).Count()
	if userPaymentNums >= accountNums {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "withdrawalAccountHasReachedTheUpperLimit"), -1)
		return
	}

	nowTime := time.Now()
	_, err = models.NewUserWalletAccount(nil).
		Field("admin_id", "user_id", "payment_id", "name", "real_name", "card_number", "address", "updated_at", "created_at").
		Args(paymentInfo.AdminId, claims.UserId, paymentInfo.Id, params.Name, params.RealName, params.CardNumber, params.Address, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
