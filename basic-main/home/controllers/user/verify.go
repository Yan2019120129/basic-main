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

type verifyParams struct {
	RealName  string `json:"real_name"`
	IdNumber  string `json:"id_number"`
	Email     string `json:"email"`
	CountryId int64  `json:"country_id"`
	Telephone string `json:"telephone"`
	Address   string `json:"address"`
	IdPhoto1  string `json:"id_photo1"`
	IdPhoto2  string `json:"id_photo2"`
	IdPhoto3  string `json:"id_photo3"`
}

// Verify 用户验证
func Verify(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(verifyParams)
	_ = body.ReadJSON(r, params)
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	claims := router.TokenManager.GetContextClaims(r)
	rds := cache.RedisPool.Get()
	defer rds.Close()
	acceptLanguage := r.Header.Get("Accept-Language")

	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)
	adminSettingList := models.NewAdminSetting(nil).GetAdminFieldAllString(settingAdminId)
	currentUserModel := models.NewUser(nil)
	currentUserModel.AndWhere("id=?", claims.UserId)
	userInfo := currentUserModel.FindOne()
	templateVerify := models.AdminSettingValueToMapInterface(adminSettingList["template_verify"])
	var realName, idNumber, email, telephone, address, idPhoto1, idPhoto2, idPhoto3 string

	// 如果需要真实姓名
	if templateVerify["real_name"].(bool) {
		if params.RealName == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "realNameCannotBeEmpty"), -1)
			return
		}
		realName = params.RealName
	}

	// 如果需要证件卡号
	if templateVerify["id_number"].(bool) {
		if params.IdNumber == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "idNumberCannotBeEmpty"), -1)
			return
		}
		idNumber = params.IdNumber
	}

	// 电子邮箱不能为空
	if templateVerify["email"].(bool) {
		if params.Email == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "emailCannotBeEmpty"), -1)
			return
		}
		email = params.Email
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
	// 手机号码不能为空
	if templateVerify["telephone"].(bool) {
		if params.Telephone == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "telephoneCannotBeEmpty"), -1)
			return
		}
		telephone = params.Telephone
	}

	// 证件地址是否为空
	if templateVerify["address"].(bool) {
		if params.Address == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "idAddressCannotBeEmpty"), -1)
			return
		}
		address = params.Address
	}

	if templateVerify["photo_front"].(bool) {
		if params.IdPhoto1 == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "idPhotoFrontCannotBeEmpty"), -1)
			return
		}
		idPhoto1 = params.IdPhoto1
	}

	if templateVerify["photo_back"].(bool) {
		if params.IdPhoto2 == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "idPhotoBackCannotBeEmpty"), -1)
			return
		}
		idPhoto2 = params.IdPhoto2
	}

	if templateVerify["photo_hold"].(bool) {
		if params.IdPhoto3 == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "holdIdPhotoCannotBeEmpty"), -1)
			return
		}
		idPhoto3 = params.IdPhoto3
	}

	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	// 更新用户信息
	if realName != "" || email != "" || telephone != "" {
		userModel := models.NewUser(tx)
		isUpdateUser := false
		if userInfo.Nickname == "" && realName != "" {
			isUpdateUser = true
			userModel.Value("nickname=?").Args(realName)
		}
		if userInfo.Email == "" && email != "" {
			isUpdateUser = true
			userModel.Value("email=?").Args(email)
		}
		if userInfo.CountryId == 0 && params.CountryId != 0 && userInfo.Telephone == "" && telephone != "" {
			isUpdateUser = true
			userModel.Value("country_id=?", "telephone=?").Args(params.CountryId, telephone)
		}

		if isUpdateUser {
			_, err = userModel.AndWhere("id=?", userInfo.Id).Update()
			if err != nil {
				panic(err)
			}
		}
	}

	verifyInfo := models.NewUserVerify(nil).AndWhere("user_id=?", userInfo.Id).FindOne()
	verifyStatus := models.UserVerifyStatusPending
	if templateVerify["autoComplete"].(bool) {
		verifyStatus = models.UserVerifyStatusComplete
	}

	nowTime := time.Now()
	if verifyInfo == nil {
		_, err = models.NewUserVerify(tx).
			Field("admin_id", "user_id", "type", "real_name", "id_number", "id_photo1", "id_photo2", "id_photo3", "address", "status", "created_at", "updated_at").
			Args(userInfo.AdminId, userInfo.Id, models.UserVerifyTypeIdCard, realName, idNumber, idPhoto1, idPhoto2, idPhoto3, address, verifyStatus, nowTime.Unix(), nowTime.Unix()).
			Insert()
		if err != nil {
			panic(err)
		}
	} else {
		_, err = models.NewUserVerify(tx).
			Value("real_name=?", "id_number=?", "id_photo1=?", "id_photo2=?", "id_photo3=?", "address=?", "data=?", "updated_at=?", "status=?").
			Args(realName, idNumber, idPhoto1, idPhoto2, idPhoto3, address, "", nowTime.Unix(), verifyStatus).
			AndWhere("id=?", verifyInfo.Id).
			Update()
		if err != nil {
			panic(err)
		}
	}

	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
