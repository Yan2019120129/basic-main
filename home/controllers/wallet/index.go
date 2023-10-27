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

type indexData struct {
	OrderSn   string  `json:"order_sn"`   //	订单号
	Name      string  `json:"name"`       //	多语言名称
	Status    int64   `json:"status"`     //	状态
	Money     float64 `json:"money"`      //	金额
	Data      string  `json:"data"`       //	数据
	CreatedAt int64   `json:"created_at"` //	时间
}

// Index 钱包记录列表
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := router.TokenManager.GetContextClaims(r)

	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)
	acceptLanguage := r.Header.Get("Accept-Language")
	rds := cache.RedisPool.Get()
	defer rds.Close()

	data := make([]*indexData, 0)
	models.NewUserWalletOrder(nil).
		Field("order_sn", "type", "money", "status", "data", "created_at").
		AndWhere("user_id=?", claims.UserId).AndWhere("status>?", models.WalletOrderStatusDelete).
		OrderBy("created_at desc").OffsetLimit(0, 20).
		Query(func(rows *sql.Rows) {
			temp := new(indexData)
			var tempType int64
			_ = rows.Scan(&temp.OrderSn, &tempType, &temp.Money, &temp.Status, &temp.Data, &temp.CreatedAt)

			switch tempType {
			case models.WalletOrderTypeDeposit:
				temp.Name = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "deposit")
			case models.WalletOrderTypeSystemDeposit:
				temp.Name = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "systemDeposit")
			case models.WalletOrderTypeWithdraw:
				temp.Name = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "withdraw")
			case models.WalletOrderTypeSystemWithdraw:
				temp.Name = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "systemDeduction")
			case models.WalletOrderTypeAssetsWithdraw:
				temp.Name = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "assetsWithdraw")
			default:
				panic("home/controllers/wallet/index.go | 没有当前钱包订单类型")
			}
			data = append(data, temp)
		})

	body.SuccessJSON(w, data)
}
