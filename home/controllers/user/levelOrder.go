package user

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type levelOrderParams struct {
	Id int64 `json:"id" validate:"required"`
}

// LevelOrder 购买等级订单
func LevelOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(levelOrderParams)
	_ = body.ReadJSON(r, params)
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	acceptLanguage := r.Header.Get("Accept-Language")
	rds := cache.RedisPool.Get()
	defer rds.Close()

	nowTime := time.Now()
	claims := router.TokenManager.GetContextClaims(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)

	userModel := models.NewUser(nil)
	userModel.AndWhere("id=?", claims.UserId)
	userInfo := userModel.FindOne()

	// 判断等级ID不否存在
	levelModel := models.NewUserLevel(nil)
	levelModel.AndWhere("id=?", params.Id).AndWhere("admin_id=?", settingAdminId)
	levelInfo := levelModel.FindOne()
	if levelInfo == nil {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "noPermissionToExecute"), -1)
		return
	}

	// 获取当前用户等级ID
	buyLevelMoney := levelInfo.Money
	userLevelOrderModel := models.NewUserLevelOrder(nil)
	userLevelOrderModel.AndWhere("user_id=?", userInfo.Id).AndWhere("status=?", models.UserLevelStatusActivate).AndWhere("updated_at>?", nowTime.Unix())
	userLevelOrderInfo := userLevelOrderModel.FindOne()

	//	如果当前已经购买过会员, 那么升级会员
	billType := models.UserBillTypeBuyLevel
	if userLevelOrderInfo != nil {
		billType = models.UserBillTypeBuyUpgradeLevel
		userLevelModel := models.NewUserLevel(nil)
		userLevelModel.AndWhere("id=?", userLevelOrderInfo.LevelId)
		userLevelInfo := userLevelModel.FindOne()

		if userLevelInfo == nil || userLevelInfo.Level >= levelInfo.Level {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "purchaseLevelFailed"), -1)
			return
		}

		//	如果是补价模式
		buyLevelMode := models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "buy_level_mode")
		if buyLevelMode == models.AdminSettingBuyLevelModePremium && buyLevelMoney > userLevelInfo.Money {
			buyLevelMoney = buyLevelMoney - userLevelInfo.Money
		}
	}

	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	updatedAt := nowTime.Unix() + 5000*31104000
	if levelInfo.Days > 0 {
		updatedAt = nowTime.Unix() + 86400*levelInfo.Days
	}
	levelId, err := models.NewUserLevelOrder(tx).
		Field("admin_id", "user_id", "level_id", "created_at", "updated_at").
		Args(userInfo.AdminId, userInfo.Id, levelInfo.Id, nowTime.Unix(), updatedAt).
		Insert()
	if err != nil {
		panic(err)
	}

	err = models.UserSpend(tx, userInfo, levelId, billType, buyLevelMoney)
	if err != nil {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, err.Error()), -1)
		return
	}

	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
