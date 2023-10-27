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

type withdrawParams struct {
	AccountId   int64   `json:"account_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	SecurityKey string  `json:"security_key"`
}

// Withdraw 提现
func Withdraw(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(withdrawParams)
	_ = body.ReadJSON(r, params)
	claims := router.TokenManager.GetContextClaims(r)
	nowTime := time.Now()

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)
	adminSettingList := models.NewAdminSetting(nil).GetAdminFieldAllString(settingAdminId)
	acceptLanguage := r.Header.Get("Accept-Language")
	rds := cache.RedisPool.Get()
	defer rds.Close()

	walletAccountModel := models.NewUserWalletAccount(nil)
	walletAccountModel.AndWhere("user_id=?", claims.UserId).AndWhere("id=?", params.AccountId).AndWhere("status=?", models.UserWalletAccountStatusActivate)
	walletAccountInfo := walletAccountModel.FindOne()
	if walletAccountInfo == nil {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "noPermissionToExecute"), -1)
		return
	}

	userModel := models.NewUser(nil)
	userModel.AndWhere("id=?", claims.UserId)
	userInfo := userModel.FindOne()
	templateWallet := models.AdminSettingValueToMapInterface(adminSettingList["template_wallet"])
	templateFreeze := models.AdminSettingValueToMapInterface(adminSettingList["template_freeze"])

	// 判断用户是否冻结状态
	if templateFreeze["withdraw"].(bool) && userInfo.Status == models.UserStatusFreeze {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "accountIsFrozen"), -1)
		return
	}

	// 如果当前用户有未完成订单，那么提示请联系客服人员
	unfulfilledOrderModel := models.NewUserWalletOrder(nil)
	unfulfilledOrderModel.AndWhere("user_id=?", claims.UserId).AndWhere("type=?", models.WalletOrderTypeWithdraw).AndWhere("status=?", models.WalletOrderStatusPending)
	unfulfilledOrderInfo := unfulfilledOrderModel.FindOne()
	if unfulfilledOrderInfo != nil {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "thereAreStillOutstandingOrders"), -1)
		return
	}

	// 判断是否提现范围
	withdrawRangeList := models.AdminSettingValueToMapInterface(adminSettingList["finance_withdraw_range"])
	withdrawMinAmount := withdrawRangeList["min"].(float64)
	withdrawMaxAmount := withdrawRangeList["max"].(float64)
	if params.Amount < withdrawMinAmount || params.Amount > withdrawMaxAmount {
		errMsg := locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "withdrawRangeError")
		errMsg = strings.ReplaceAll(errMsg, "{min}", strconv.FormatFloat(withdrawMinAmount, 'f', 0, 64))
		errMsg = strings.ReplaceAll(errMsg, "{max}", strconv.FormatFloat(withdrawMaxAmount, 'f', 10, 64))
		body.ErrorJSON(w, errMsg, -1)
		return
	}

	//	获取提现时间范围, [{"staTime": "09:00:00", "endTime": "12:00:00"}, {"staTime": "12:00:00", "endTime": "18:00:00"}]
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))
	withdrawTimeList := models.AdminSettingValueToMapInterfaces(adminSettingList["finance_withdraw_times"])
	withdrawNowTime := time.Now().In(location)
	if len(withdrawTimeList) > 0 {
		isTimeRange := false
		timeRangeError := locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "withdrawTimeError")
		for i := 0; i < len(withdrawTimeList); i++ {
			staTime, _ := time.ParseInLocation("2006-01-02 15:04:05", withdrawNowTime.Format("2006-01-02")+" "+withdrawTimeList[i]["sta_time"].(string), location)
			endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", withdrawNowTime.Format("2006-01-02")+" "+withdrawTimeList[i]["end_time"].(string), location)

			// 如果结束时间比较小，那么+1天时间
			minTimeUnix := staTime.Unix()
			maxTimeUnix := endTime.Unix()
			if maxTimeUnix < minTimeUnix {
				maxTimeUnix += 86400
			}

			if withdrawNowTime.Unix() > minTimeUnix && withdrawNowTime.Unix() < maxTimeUnix {
				isTimeRange = true
			}
			timeRangeError = strings.ReplaceAll(timeRangeError, "{timeRange}", withdrawTimeList[i]["sta_time"].(string)+" - "+withdrawTimeList[i]["end_time"].(string))
		}

		if !isTimeRange {
			body.ErrorJSON(w, timeRangeError, -1)
			return
		}
	}

	// 判断今日是否超过提现次数
	settingWithdrawNums := models.AdminSettingValueToMapInterface(adminSettingList["finance_withdraw_nums"])
	beforeTime := time.Now().AddDate(0, 0, int(settingWithdrawNums["days"].(float64)))
	userTodayWithdrawNums := models.NewUserWalletOrder(nil).
		AndWhere("user_id=?", userInfo.Id).
		AndWhere("type=?", models.WalletOrderTypeWithdraw).
		AndWhere("status=?", models.WalletOrderStatusComplete).
		AndWhere("updated_at>?", beforeTime.Unix()).
		Count()

	if userTodayWithdrawNums >= int64(settingWithdrawNums["nums"].(float64)) {
		errMsg := locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "withdrawNumsError")
		withdrawDaysStr := strconv.FormatInt(int64(settingWithdrawNums["days"].(float64)), 10)
		withdrawNumsStr := strconv.FormatInt(int64(settingWithdrawNums["nums"].(float64)), 10)
		errMsg = strings.ReplaceAll(errMsg, "{days}", withdrawDaysStr)
		errMsg = strings.ReplaceAll(errMsg, "{nums}", withdrawNumsStr)

		body.ErrorJSON(w, errMsg, -1)
		return
	}

	// 判断是否余额不足
	if userInfo.Money < params.Amount {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "insufficientBalance"), -1)
		return
	}

	// 判断是否开启安全密钥，判断是否正确
	if templateWallet["withdraw_security_key"].(bool) {
		if params.SecurityKey == "" || utils.PasswordEncrypt(params.SecurityKey) != userInfo.SecurityKey {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "incorrectSecurityKey"), -1)
			return
		}
	}

	// 如果没有实名，不能提现
	userVerifyInfo := models.NewUserVerify(nil).AndWhere("user_id=?", userInfo.Id).FindOne()
	if templateWallet["withdraw_verify"].(bool) && userVerifyInfo != nil && userVerifyInfo.Status != models.UserVerifyStatusComplete {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "notYetRegistered"), -1)
		return
	}

	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	// 写入订单数据库
	var withdrawOrderFee float64
	if adminSettingList["finance_withdraw_fee"] != "" {
		withdrawOrderFee, _ = strconv.ParseFloat(adminSettingList["finance_withdraw_fee"], 64)
	}

	walletOrderId, err := models.NewUserWalletOrder(tx).
		Field("order_sn", "admin_id", "user_id", "type", "payment_id", "money", "fee", "data", "balance", "proof", "updated_at", "created_at").
		Args(utils.NewRandom().OrderSn(), userInfo.AdminId, userInfo.Id, models.WalletOrderTypeWithdraw, walletAccountInfo.Id, params.Amount, withdrawOrderFee, "", userInfo.Money, "", nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	// 写入数据
	err = models.UserSpend(tx, userInfo, walletOrderId, models.UserBillTypeWithdraw, params.Amount)
	if err != nil {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, err.Error()), -1)
		return
	}

	// 提现提示音
	//	TODO...
	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
