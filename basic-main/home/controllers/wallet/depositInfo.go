package wallet

import (
	"basic/models"
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type depositInfoItem struct {
	Id          int64  `json:"id"`           //	ID
	Icon        string `json:"icon"`         //	图标
	Type        int64  `json:"type"`         //	1银行转账 10数字货币 20三方支付
	Name        string `json:"name"`         //	建设银行 波场链
	AccountName string `json:"account_name"` //	真实姓名｜USDT
	AccountCode string `json:"account_code"` //	卡号｜地址
	Data        string `json:"data"`         //	额外数据
}

type depositInfoData struct {
	Tips  string             `json:"tips"`
	Items []*depositInfoItem `json:"items"`
}

// DepositInfo 充值列表
func DepositInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := router.TokenManager.GetContextClaims(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)
	acceptLanguage := r.Header.Get("Accept-Language")
	rds := cache.RedisPool.Get()
	defer rds.Close()

	data := make([]*depositInfoItem, 0)
	models.NewWalletPayment(nil).Field("id", "icon", "type", "name", "account_name", "account_code", "data").
		AndWhere("admin_id=?", settingAdminId).AndWhere("status=?", models.WalletPaymentStatusActivate).AndWhere("mode=?", models.WalletPaymentModeDeposit).
		OrderBy("sort asc").
		Query(func(rows *sql.Rows) {
			temp := new(depositInfoItem)
			_ = rows.Scan(&temp.Id, &temp.Icon, &temp.Type, &temp.Name, &temp.AccountName, &temp.AccountCode, &temp.Data)
			data = append(data, temp)
		})

	tips := models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "finance_deposit_tip")
	body.SuccessJSON(w, &depositInfoData{
		Items: data,
		Tips:  locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, tips),
	})
}
