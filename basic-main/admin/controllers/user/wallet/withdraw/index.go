package withdraw

import (
	"basic/models"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type indexParams struct {
	AdminName         string                 `json:"admin_name"`
	UserName          string                 `json:"username"`
	OrderSn           string                 `json:"order_sn"`
	PaymentRealName   string                 `json:"payment_real_name"`
	PaymentCardNumber string                 `json:"payment_card_number"`
	Status            int64                  `json:"status"`
	Data              string                 `json:"data"`
	DateTime          *define.RangeTimeParam `json:"created_at"`
	Pagination        *define.Pagination     `json:"pagination"` //	分页
}

type indexData struct {
	models.UserWalletOrderAttrs
	MoneyFee          float64 `json:"money_fee"`
	AdminName         string  `json:"admin_name"`
	UserName          string  `json:"username"`
	PaymentName       string  `json:"payment_name"`
	PaymentRealName   string  `json:"payment_real_name"`
	PaymentCardNumber string  `json:"payment_card_number"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取子级包括自己ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	model := models.NewUserWalletOrder(nil)
	model.AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("type>=?", models.WalletOrderTypeWithdraw).AndWhere("status<>?", models.WalletOrderStatusDelete)

	define.NewFilterEmpty(model.Db).
		String("order_sn=?", params.OrderSn).
		String("data=?", params.Data).
		Int64("status=?", params.Status).
		RangeTime("created_at between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)
	// 管理员名称
	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	// 用户名称
	if params.UserName != "" {
		model.Db.AndWhere("user_id in (" + strings.Join(models.NewUser(nil).FindUserLikeNameIds(params.UserName), ",") + ")")
	}

	// 支付方式 姓名｜Token
	if params.PaymentRealName != "" {
		paymentRealNameIds := models.NewUserWalletAccount(nil).Field("id").AndWhere("real_name like ?", "%"+params.PaymentRealName+"%").ColumnString()
		if len(paymentRealNameIds) == 0 {
			paymentRealNameIds = append(paymentRealNameIds, "-1")
		}
		model.Db.AndWhere("payment_id in (" + strings.Join(paymentRealNameIds, ",") + ")")
	}

	// 支付方式 卡号｜地址
	if params.PaymentCardNumber != "" {
		paymentCardNumberIds := models.NewUserWalletAccount(nil).Field("id").AndWhere("card_number like ?", "%"+params.PaymentCardNumber+"%").ColumnString()
		if len(paymentCardNumberIds) == 0 {
			paymentCardNumberIds = append(paymentCardNumberIds, "-1")
		}
		model.Db.AndWhere("payment_id in (" + strings.Join(paymentCardNumberIds, ",") + ")")
	}

	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.OrderSn, &tmp.AdminId, &tmp.UserId, &tmp.UserType, &tmp.Type, &tmp.PaymentId, &tmp.Money, &tmp.Balance, &tmp.Status, &tmp.Proof, &tmp.Data, &tmp.Fee, &tmp.UpdatedAt, &tmp.CreatedAt)
		// 当前用户信息
		userInfo := models.NewUser(nil).AndWhere("id=?", tmp.UserId).FindOne()
		if userInfo != nil {
			tmp.UserName = userInfo.UserName
		}

		// 当前管理员信息
		adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}

		// 支付名称
		paymentInfo := models.NewUserWalletAccount(nil).AndWhere("id=?", tmp.PaymentId).FindOne()
		if paymentInfo != nil {
			tmp.PaymentName = paymentInfo.Name
			tmp.PaymentRealName = paymentInfo.RealName
			tmp.PaymentCardNumber = paymentInfo.CardNumber
		}

		// 扣除手续费金额
		tmp.MoneyFee = tmp.Money - (tmp.Money * tmp.Fee / 100)
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
