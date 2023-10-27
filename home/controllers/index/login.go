package index

import (
	"basic/models"
	"database/sql"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type loginParams struct {
	UserName     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	CaptchaId    string `json:"captcha_id"`
	CaptchaValue string `json:"captcha_value"`
}

type LoginData struct {
	Info     *models.UserInfo `json:"info"`     //	用户信息
	TokenKey string           `json:"tokenKey"` //	TokenKey
	Token    string           `json:"token"`    //	Token
}

// Login 用户登陆
func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(loginParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	adminId := models.NewAdminUser(nil).GetDomainAdminId(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	acceptLanguage := r.Header.Get("Accept-Language")
	settingAdminList := models.NewAdminSetting(nil).GetAdminFieldAllString(settingAdminId)
	templateLogin := models.AdminSettingValueToMapInterface(settingAdminList["template_login"])
	nowTime := time.Now()
	rds := cache.RedisPool.Get()
	defer rds.Close()

	// 验证码是否正确
	if templateLogin["show_code"].(bool) {
		if !captcha.VerifyString(params.CaptchaId, params.CaptchaValue) {
			body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "wrongVerificationCode"), -1)
			return
		}
	}

	//	获取用户信息 [是否存在，密码是否正确]
	userModel := models.NewUser(nil)
	userModel.AndWhere("username=?", params.UserName).AndWhere("status=?", models.UserStatusActivate).AndWhere("type>=?", models.UserTypeVirtual)
	userInfo := userModel.FindOne()
	if userInfo == nil || userInfo.Password != utils.PasswordEncrypt(params.Password) {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "incorrectAccountOrPassword"), -1)
		return
	}

	//	登陆成功修改更新时间｜IP4地址
	_, err = models.NewUser(nil).Value("updated_at=?", "ip4=INET_ATON(?)").
		Args(nowTime.Unix(), utils.GetUserRealIP(r)).
		AndWhere("id=?", userInfo.Id).Update()
	if err != nil {
		panic(err)
	}

	//	验证信息
	verifyInfo := new(models.UserVerifyInfo)
	models.NewUserVerify(nil).Field("status", "data").AndWhere("user_id=?", userInfo.Id).QueryRow(func(row *sql.Row) {
		_ = row.Scan(&verifyInfo.Status, &verifyInfo.Data)
	})

	//	等级信息
	levelInfo := new(models.UserLevelInfo)
	models.NewUserLevelOrder(nil).Field("level_id").AndWhere("user_id=?", userInfo.Id).AndWhere("status=?", models.UserLevelOrderStatusActivate).QueryRow(func(row *sql.Row) {
		var userLevelId int64
		_ = row.Scan(&userLevelId)

		if userLevelId > 0 {
			models.NewUserLevel(nil).Field("id", "level", "name", "icon", "days", "money", "created_at", "updatedAt").AndWhere("id=?", userLevelId).QueryRow(func(row *sql.Row) {
				_ = row.Scan(&levelInfo.Id, &levelInfo.Level, &levelInfo.Name, &levelInfo.Icon, &levelInfo.Days, &levelInfo.Money, &levelInfo.CreatedAt, &levelInfo.UpdatedAt)
			})
		}
		if levelInfo.Name != "" {
			levelInfo.Name = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, levelInfo.Name)
		}
	})

	tokenKey := models.TokenParamsPrefix(models.HomePrefixTokenKey, settingAdminId)
	body.SuccessJSON(w, &LoginData{
		Info: &models.UserInfo{
			Id:          userInfo.Id,
			CountryId:   userInfo.CountryId,
			UserName:    userInfo.UserName,
			Nickname:    userInfo.Nickname,
			Email:       userInfo.Email,
			Telephone:   userInfo.Telephone,
			Avatar:      userInfo.Avatar,
			Sex:         userInfo.Sex,
			Birthday:    userInfo.Birthday,
			Money:       userInfo.Money,
			FreezeMoney: userInfo.FreezeMoney,
			Data:        userInfo.Data,
			VerifyInfo:  verifyInfo,
			LevelInfo:   levelInfo,
			InviteCode:  models.NewUserInvite(nil).GetInviteCode(userInfo.AdminId, userInfo.Id),
			CreatedAt:   userInfo.CreatedAt,
			UpdatedAt:   userInfo.UpdatedAt,
		},
		TokenKey: tokenKey,
		Token:    router.TokenManager.Generate(rds, tokenKey, userInfo.AdminId, userInfo.Id),
	})
}
