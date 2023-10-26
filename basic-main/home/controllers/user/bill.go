package user

import (
	"basic/models"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type billData struct {
	SumAmount float64       `json:"sumAmount"` //	合计
	Items     []*billItem   `json:"items"`     //	数据
	Options   []*billOption `json:"options"`   //	账单类型数组
}

type billOption struct {
	Label string `json:"label"` //	账单类型名称
	Value int64  `json:"value"` //	账单类型ID
}

type billItem struct {
	Name           string  `json:"name"`            //	多语言键名
	SourceUserName string  `json:"source_username"` //	来源名称
	Type           int64   `json:"type"`            //	账单类型
	Money          float64 `json:"money"`           //	账单金额
	Data           string  `json:"data"`            //	数据
	CreatedAt      int64   `json:"created_at"`      //	创建时间
}

type billParams struct {
	Ids      []string               `json:"ids"`  //	类型ID组
	DateTime *define.RangeTimeParam `json:"date"` //	时间范围
}

// Bill 用户账单
func Bill(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(billParams)
	_ = body.ReadJSON(r, params)

	claims := router.TokenManager.GetContextClaims(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)
	rds := cache.RedisPool.Get()
	defer rds.Close()

	//	账单类型名称
	options := []*billOption{
		{Label: locales.Manager.GetAdminLocales(rds, settingAdminId, r.Header.Get("Accept-Language"), models.UserBillTypeNameMap[models.UserBillTypeSystemDeposit]), Value: models.UserBillTypeSystemDeposit},
		{Label: locales.Manager.GetAdminLocales(rds, settingAdminId, r.Header.Get("Accept-Language"), models.UserBillTypeNameMap[models.UserBillTypeSystemDeduction]), Value: models.UserBillTypeSystemDeduction},
		{Label: locales.Manager.GetAdminLocales(rds, settingAdminId, r.Header.Get("Accept-Language"), models.UserBillTypeNameMap[models.UserBillTypeDeposit]), Value: models.UserBillTypeDeposit},
		{Label: locales.Manager.GetAdminLocales(rds, settingAdminId, r.Header.Get("Accept-Language"), models.UserBillTypeNameMap[models.UserBillTypeWithdraw]), Value: models.UserBillTypeWithdraw},
		{Label: locales.Manager.GetAdminLocales(rds, settingAdminId, r.Header.Get("Accept-Language"), models.UserBillTypeNameMap[models.UserBillTypeWithdrawRefuse]), Value: models.UserBillTypeWithdrawRefuse},
		{Label: locales.Manager.GetAdminLocales(rds, settingAdminId, r.Header.Get("Accept-Language"), models.UserBillTypeNameMap[models.UserBillTypeRegisterRewards]), Value: models.UserBillTypeRegisterRewards},
		{Label: locales.Manager.GetAdminLocales(rds, settingAdminId, r.Header.Get("Accept-Language"), models.UserBillTypeNameMap[models.UserBillTypeInviteRewards]), Value: models.UserBillTypeInviteRewards},
		{Label: locales.Manager.GetAdminLocales(rds, settingAdminId, r.Header.Get("Accept-Language"), models.UserBillTypeNameMap[models.UserBillTypeBuyProduct]), Value: models.UserBillTypeBuyProduct},
		{Label: locales.Manager.GetAdminLocales(rds, settingAdminId, r.Header.Get("Accept-Language"), models.UserBillTypeNameMap[models.UserBillTypeProductProfit]), Value: models.UserBillTypeProductProfit},
	}

	var sumAmount float64 = 0
	data := make([]*billItem, 0)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	model := models.NewUserBill(nil)
	model.Field("name", "type", "money", "data", "created_at").
		AndWhere("user_id=?", claims.UserId)
	define.NewFilterEmpty(model.Db).
		RangeTime("created_at between ? and ?", params.DateTime, location)
	if len(params.Ids) > 0 {
		model.AndWhere("type in (" + strings.Join(params.Ids, ",") + ")")
	}

	model.OrderBy("id desc").OffsetLimit(0, 50).
		Query(func(rows *sql.Rows) {
			temp := new(billItem)
			_ = rows.Scan(&temp.Name, &temp.Type, &temp.Money, &temp.Data, &temp.CreatedAt)
			//	自动计算类型金额是否正负金额
			temp.Money = models.GetBillTypeMoney(temp.Type, 0, temp.Money)

			sumAmount += temp.Money
			temp.Name = locales.Manager.GetAdminLocales(rds, settingAdminId, r.Header.Get("Accept-Language"), temp.Name)
			data = append(data, temp)
		})

	body.SuccessJSON(w, &billData{Items: data, SumAmount: sumAmount, Options: options})
}
