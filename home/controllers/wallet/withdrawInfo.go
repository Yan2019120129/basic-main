package wallet

import (
	"basic/models"
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type withdrawInfoItem struct {
	Id          int64  `json:"id"`           //	ID
	Icon        string `json:"icon"`         //	图标
	Type        int64  `json:"type"`         //	1银行转账 10数字货币 20三方支付
	Name        string `json:"name"`         //	建设银行 波场链
	AccountName string `json:"account_name"` //	默认名字
	Data        string `json:"data"`         //	额外数据
}

// WithdrawInfo 提现信息
func WithdrawInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := router.TokenManager.GetContextClaims(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)

	data := make([]*withdrawInfoItem, 0)
	models.NewWalletPayment(nil).Field("id", "icon", "type", "name", "account_name", "data").
		AndWhere("admin_id=?", settingAdminId).AndWhere("status=?", models.WalletPaymentStatusActivate).
		AndWhere("mode=?", models.WalletPaymentModeWithdraw).
		OrderBy("sort asc").
		Query(func(rows *sql.Rows) {
			temp := new(withdrawInfoItem)
			_ = rows.Scan(&temp.Id, &temp.Icon, &temp.Type, &temp.Name, &temp.AccountName, &temp.Data)
			data = append(data, temp)
		})

	body.SuccessJSON(w, data)
}
