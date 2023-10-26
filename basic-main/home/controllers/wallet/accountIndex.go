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

type accountIndexItem struct {
	PaymentId          int64  `json:"payment_id"`
	PaymentIcon        string `json:"payment_icon"`
	PaymentType        int64  `json:"payment_type"`
	PaymentName        string `json:"payment_name"`
	PaymentAccountName string `json:"payment_account_name"`
	Id                 int64  `json:"id"`
	Name               string `json:"name"`
	RealName           string `json:"real_name"`
	CardNumber         string `json:"card_number"`
	Address            string `json:"address"`
}

type accountIndexData struct {
	Items []*accountIndexItem `json:"items"`
	Tips  string              `json:"tips"`
}

// AccountIndex 用户钱包账户列表
func AccountIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := router.TokenManager.GetContextClaims(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)
	acceptLanguage := r.Header.Get("Accept-Language")
	rds := cache.RedisPool.Get()
	defer rds.Close()

	data := make([]*accountIndexItem, 0)
	models.NewUserWalletAccount(nil).Field("id", "payment_id", "name", "real_name", "card_number", "address").
		AndWhere("user_id=?", claims.UserId).AndWhere("status=?", models.UserWalletAccountStatusActivate).
		OrderBy("sort asc").OffsetLimit(0, 10).
		Query(func(rows *sql.Rows) {
			temp := new(accountIndexItem)
			_ = rows.Scan(&temp.Id, &temp.PaymentId, &temp.Name, &temp.RealName, &temp.CardNumber, &temp.Address)

			paymentInfo := models.NewWalletPayment(nil).AndWhere("id=?", temp.PaymentId).FindOne()
			temp.PaymentName = paymentInfo.Name
			temp.PaymentIcon = paymentInfo.Icon
			temp.PaymentType = paymentInfo.Type
			temp.PaymentAccountName = paymentInfo.AccountName

			data = append(data, temp)
		})

	tips := models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "finance_withdraw_tip")
	body.SuccessJSON(w, &accountIndexData{
		Items: data,
		Tips:  locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, tips),
	})
}
