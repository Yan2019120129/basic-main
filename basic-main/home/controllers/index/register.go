package index

import (
	"basic/models"
	"net/http"
	"strconv"
	"time"

	"github.com/dchest/captcha"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type registerParams struct {
	UserName     string `json:"username" validate:"required"` // 用户名
	Password     string `json:"password" validate:"required"` // 用户密码
	SecurityKey  string `json:"security_key"`                 // 安全密钥
	Nickname     string `json:"nickname"`                     // 昵称
	Email        string `json:"email"`                        // 邮箱
	CountryId    int64  `json:"country_id"`                   // 国家ID
	Telephone    string `json:"telephone"`                    // 手机号码
	InviteCode   string `json:"invite_code"`                  // 邀请码
	CaptchaId    string `json:"captcha_id"`                   // 验证码ID
	CaptchaValue string `json:"captcha_value"`                // 验证码
}

// Register 用户注册
func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(registerParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	//	判断管理员是否存在
	adminId := models.NewAdminUser(nil).GetDomainAdminId(r)
	adminInfo := models.NewAdminUser(nil).AndWhere("id=?", adminId).FindOne()
	if adminInfo == nil {
		panic("home/controllers/index/register.go | 注册管理域名不存在")
	}

	rds := cache.RedisPool.Get()
	defer rds.Close()

	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	var parentId, countryId int64
	var nickname, email, securityKey, telephone string
	userIP := utils.GetUserRealIP(r)
	acceptLanguage := r.Header.Get("Accept-Language")
	settingAdminList := models.NewAdminSetting(nil).GetAdminFieldAllString(settingAdminId)
	templateRegister := models.AdminSettingValueToMapInterface(settingAdminList["template_register"])

	// 验证码是否正确
	if templateRegister["show_code"].(bool) {
		if !captcha.VerifyString(params.CaptchaId, params.CaptchaValue) {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "wrongVerificationCode"), -1)
			return
		}
	}

	// 判断用户名是否存在
	userModel := models.NewUser(nil)
	userModel.AndWhere("username=?", params.UserName)
	userInfo := userModel.FindOne()
	if userInfo != nil {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "userNameAlreadyExists"), -1)
		return
	}

	// 如果设置了 "昵称" "邮箱" "手机号码" "安全密钥" "邀请码"
	if templateRegister["nickname"].(bool) {
		if params.Nickname == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "nickNameCannotBeEmpty"), -1)
			return
		}
		nickname = params.Nickname
	}

	// 是否需要邮箱
	if templateRegister["email"].(bool) {
		if params.Email == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "emailCannotBeEmpty"), -1)
			return
		}
		email = params.Email
	}

	// 是否需要安全密钥
	if templateRegister["security_key"].(bool) {
		if params.SecurityKey == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "securityKeyCannotBeEmpty"), -1)
			return
		}
		securityKey = utils.PasswordEncrypt(params.SecurityKey)
	}

	// 是否需要邀请码
	if templateRegister["invite_code"].(bool) && params.InviteCode == "" {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "inviteCodeCannotBeEmpty"), -1)
		return
	}

	if params.InviteCode != "" {
		inviteInfo := models.NewUserInvite(nil).AndWhere("code=?", params.InviteCode).FindOne()
		if inviteInfo == nil {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "invitationCodeIsIncorrect"), -1)
			return
		}
		adminId = inviteInfo.AdminId
		parentId = inviteInfo.UserId
	}

	// 是否需要手机号码
	if templateRegister["telephone"].(bool) {
		if params.Telephone == "" {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "telephoneCannotBeEmpty"), -1)
			return
		}

		userPhoneModel := models.NewUser(nil)
		userPhoneModel.AndWhere("telephone=?", params.Telephone)
		userPhoneInfo := userPhoneModel.FindOne()
		if userPhoneInfo != nil {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "telephoneAlreadyExists"), -1)
			return
		}

		countryModel := models.NewCountry(nil)
		countryModel.AndWhere("id=?", params.CountryId).AndWhere("status=?", models.CountryStatusActivate)
		countryInfo := countryModel.FindOne()
		if countryInfo == nil {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "telephoneAreaCodeDoesNotExist"), -1)
			return
		}
		telephone = params.Telephone
		countryId = countryInfo.Id
	} else {
		//	如果不需要手机号码,那么自动获取
		countryId, _ = models.NewCountry(nil).AutoUserCountry(settingAdminId, userIP)
	}

	//	如果安全密钥没有设置, 那么设置跟密码一样
	if securityKey == "" {
		securityKey = utils.PasswordEncrypt(params.Password)
	}

	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	nowTime := time.Now()
	userId, err := models.NewUser(tx).
		Field("admin_id", "ip4", "parent_id", "country_id", "username", "money", "nickname", "email", "telephone", "password", "security_key", "created_at", "updated_at").
		Value("?", "INET_ATON(?)", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?").
		Args(adminId, userIP, parentId, countryId, params.UserName, 0, nickname, email, telephone, utils.PasswordEncrypt(params.Password), securityKey, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	//	邀请奖励
	parentUserRewards, _ := strconv.ParseFloat(settingAdminList["invite_rewards"], 64)
	if parentId > 0 && parentUserRewards > 0 {
		parentUserModel := models.NewUser(nil)
		parentUserModel.AndWhere("id=?", parentId).AndWhere("status=?", models.UserStatusActivate)
		parentUserInfo := parentUserModel.FindOne()

		if parentUserInfo != nil {
			_ = models.UserDeposit(tx, parentUserInfo, userId, models.UserBillTypeInviteRewards, parentUserRewards)
		}
	}

	// 	注册奖励
	inviteAmount, _ := strconv.ParseFloat(settingAdminList["register_rewards"], 64)
	if inviteAmount > 0 {
		currentUserInfo := &models.UserAttrs{Id: userId, AdminId: adminId}
		_ = models.UserDeposit(tx, currentUserInfo, 0, models.UserBillTypeRegisterRewards, inviteAmount)
	}

	_ = tx.Commit()
	tokenKey := models.TokenParamsPrefix(models.HomePrefixTokenKey, settingAdminId)
	body.SuccessJSON(w, &LoginData{
		Info: &models.UserInfo{
			Id:          userId,
			CountryId:   countryId,
			UserName:    params.UserName,
			Nickname:    nickname,
			Email:       email,
			Telephone:   telephone,
			Avatar:      "",
			Sex:         0,
			Birthday:    0,
			Money:       inviteAmount,
			FreezeMoney: 0,
			Data:        "",
			VerifyInfo:  new(models.UserVerifyInfo),
			LevelInfo:   new(models.UserLevelInfo),
			InviteCode:  models.NewUserInvite(nil).GetInviteCode(adminId, userId),
			CreatedAt:   nowTime.Unix(),
			UpdatedAt:   nowTime.Unix(),
		},
		TokenKey: tokenKey,
		Token:    router.TokenManager.Generate(rds, tokenKey, adminId, userId),
	})
}
