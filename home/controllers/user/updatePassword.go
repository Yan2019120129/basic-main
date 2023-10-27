package user

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

type updatePasswordParams struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

// UpdatePassword 更新用户登陆密码
func UpdatePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updatePasswordParams)
	_ = body.ReadJSON(r, params)

	claims := router.TokenManager.GetContextClaims(r)
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	rds := cache.RedisPool.Get()
	defer rds.Close()
	acceptLanguage := r.Header.Get("Accept-Language")

	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)
	adminSettingList := models.NewAdminSetting(nil).GetAdminFieldAllString(settingAdminId)
	userModel := models.NewUser(nil)
	userModel.AndWhere("id=?", claims.UserId)
	userInfo := userModel.FindOne()

	// 是否可以修改密码
	templateBasic := models.AdminSettingValueToMapInterface(adminSettingList["template_basic"])
	if !templateBasic["update_password"].(bool) {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "noPermissionToExecute"), -1)
		return
	}

	if utils.PasswordEncrypt(params.OldPassword) != userInfo.Password {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "incorrectOldPassword"), -1)
		return
	}

	_, err = models.NewUser(nil).
		Value("password=?").
		Args(utils.PasswordEncrypt(params.NewPassword)).
		AndWhere("id=?", userInfo.Id).
		Update()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
