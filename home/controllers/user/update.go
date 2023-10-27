package user

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type updateParams struct {
	Avatar    string `json:"avatar"`
	NickName  string `json:"nickname"`
	Email     string `json:"email"`
	CountryId int64  `json:"country_id"`
	Telephone string `json:"telephone"`
	Sex       int64  `json:"sex"`
	Birthday  string `json:"birthday"`
}

// Update 更新用户信息
func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updateParams)
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
	adminSettingList := models.NewAdminSetting(nil).GetAdminFieldAllString(settingAdminId)
	templateBasic := models.AdminSettingValueToMapInterface(adminSettingList["template_basic"])

	//	判断是否开启更新权限
	if (!templateBasic["update_avatar"].(bool) && params.Avatar != "") ||
		(!templateBasic["update_email"].(bool) && params.NickName != "") ||
		(!templateBasic["update_nickname"].(bool) && params.Email != "") ||
		(!templateBasic["update_sex"].(bool) && params.Telephone != "") ||
		(!templateBasic["update_telephone"].(bool) && params.Sex != 0) ||
		(!templateBasic["update_birthday"].(bool) && params.Birthday != "") {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "noPermissionToExecute"), -1)
		return
	}

	// 判断国家是否存在
	if params.CountryId != 0 {
		countryModel := models.NewCountry(nil)
		countryModel.AndWhere("id=?", params.CountryId).AndWhere("admin_id=?", settingAdminId)
		countryInfo := countryModel.FindOne()
		if countryInfo == nil {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "telephoneAreaCodeDoesNotExist"), -1)
			return
		}
	}

	// 判断性别是否正确
	if params.Sex != -1 && params.Sex != 1 && params.Sex != 2 {
		params.Sex = -1
	}

	model := models.NewUser(nil)
	nowTime := time.Now()
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("avatar=?", params.Avatar).
		String("nickname=?", params.NickName).
		String("email=?", params.Email).
		Int64("country_id=?", params.CountryId).
		String("telephone=?", params.Telephone).
		Int64("sex=?", params.Sex).
		DateTime("birthday=?", params.Birthday, location).
		Int64("updated_at=?", nowTime.Unix())

	_, err = model.AndWhere("id=?", claims.UserId).Update()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
