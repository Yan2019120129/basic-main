package wallet

import (
	"basic/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type accountDelete struct {
	Id          int64  `json:"id" validate:"required"`
	SecurityKey string `json:"security_key"`
}

func AccountDelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(accountDelete)
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

	//	钱包配置
	templateWallet := models.AdminSettingValueToMapInterface(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "template_wallet"))

	// 判断是否开启可以删除
	if !templateWallet["delete"].(bool) {
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
	userAccountModel := models.NewUserWalletAccount(nil)
	userAccountModel.AndWhere("id=?", params.Id).AndWhere("user_id=?", userInfo.Id)
	userAccountInfo := userAccountModel.FindOne()
	if userAccountInfo == nil {
		panic("home/controllers/wallet/accountDelete.go | 用户钱包账户异常")
	}

	_, err = models.NewUserWalletAccount(nil).Value("status=?").Args(models.UserWalletAccountStatusDelete).AndWhere("id=?", userAccountInfo.Id).Update()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
