package wallet

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type accountUpdateParams struct {
	Id          int64  `json:"id" validate:"required"`
	SecurityKey string `json:"security_key"`
	Name        string `json:"name"`        //	建设银行｜波场链
	RealName    string `json:"real_name"`   //	真实姓名｜USDT
	CardNumber  string `json:"card_number"` //	卡号｜地址
	Address     string `json:"address"`     //	银行地址
}

// AccountUpdate 用户钱包账户更新
func AccountUpdate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(accountUpdateParams)
	_ = body.ReadJSON(r, params)
	claims := router.TokenManager.GetContextClaims(r)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	acceptLanguage := r.Header.Get("Accept-Language")
	rds := cache.RedisPool.Get()
	defer rds.Close()

	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)
	userModel := models.NewUser(nil)
	userModel.AndWhere("id=?", claims.UserId)
	userInfo := userModel.FindOne()

	templateWallet := models.AdminSettingValueToMapInterface(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "template_wallet"))

	// 判断是否开启可以删除
	if !templateWallet["update"].(bool) {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "noPermissionToExecute"), -1)
		return
	}

	// 判断是否开启安全密钥，判断是否正确
	if templateWallet["security_key"].(bool) {
		if params.SecurityKey == "" || utils.PasswordEncrypt(params.SecurityKey) != userInfo.SecurityKey {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "incorrectSecurityKey"), -1)
			return
		}
	}
	// 判断当前用户是否拥有这个钱包权限

	// 判断当前用户是否拥有这个钱包权限
	userAccountModel := models.NewUserWalletAccount(nil)
	userAccountModel.AndWhere("id=?", params.Id).AndWhere("user_id=?", userInfo.Id)
	userAccountInfo := userAccountModel.FindOne()
	if userAccountInfo == nil {
		panic("home/controllers/wallet/accountUpdate.go | 用户钱包账户异常")
	}

	model := models.NewUserWalletAccount(nil)
	nowTime := time.Now()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("name=?", params.Name).
		String("real_name=?", params.RealName).
		String("card_number=?", params.CardNumber).
		String("address=?", params.Address).
		Int64("updated_at=?", nowTime.Unix())

	_, err = model.AndWhere("id=?", userAccountInfo.Id).Update()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
