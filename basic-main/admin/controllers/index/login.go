package index

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type loginParams struct {
	UserName     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	CaptchaId    string `json:"captcha_id" validate:"required"`
	CaptchaValue string `json:"captcha_value" validate:"required"`
}

type userInfo struct {
	Id         int64   `json:"id"`          //	管理员ID
	Username   string  `json:"username"`    //	用户名
	Nickname   string  `json:"nickname"`    //	昵称
	Email      string  `json:"email"`       //	邮箱
	Avatar     string  `json:"avatar"`      //	头像
	Money      float64 `json:"money"`       //	金额
	Data       string  `json:"data"`        //	数据
	InviteCode string  `json:"invite_code"` //	邀请码
	ExpiredAt  int64   `json:"expired_at"`  //	过期时间
	UpdatedAt  int64   `json:"updatedAt"`   //	更新时间
}

type loginData struct {
	Menu        []*models.AdminMenuList `json:"menu"`        //	菜单
	RouterList  []string                `json:"router_list"` //	路由列表
	UserInfo    *userInfo               `json:"info"`        //	用户信息
	Token       string                  `json:"token"`       //	Token
	TokenKey    string                  `json:"token_key"`   // 	TokenKey
	OnlineToken string                  `json:"onlineToken"` //	客服Token
	OnlineKey   string                  `json:"onlineKey"`   //	客服Key
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(loginParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 验证码是否正确
	if !captcha.VerifyString(params.CaptchaId, params.CaptchaValue) {
		body.ErrorJSON(w, "验证码错误", -1)
		return
	}

	//	获取用户信息 [是否存在，密码是否正确，是否过期]
	nowTime := time.Now().Unix()
	adminInfo := models.NewAdminUser(nil).AndWhere("username=?", params.UserName).AndWhere("status=?", models.AdminUserStatusActivate).FindOne()
	if adminInfo == nil || adminInfo.Password != utils.PasswordEncrypt(params.Password) {
		body.ErrorJSON(w, "账号或密码错误", -1)
		return
	}

	//	如果代理管理员过期, 那么也是现实已经过期
	if adminInfo.ExpiredAt > 0 && adminInfo.ExpiredAt < nowTime {
		body.ErrorJSON(w, "账号已过期, 请联系管理员", -1)
		return
	}

	// 组长管理过期时间
	if adminInfo.ParentId > 0 && adminInfo.ParentId != models.AdminUserSupermanId {
		settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminInfo.Id)
		settingAdminInfo := models.NewAdminUser(nil).AndWhere("id=?", settingAdminId).FindOne()

		if settingAdminInfo.ExpiredAt > 0 && settingAdminInfo.ExpiredAt < nowTime {
			body.ErrorJSON(w, "账号已过期, 请联系管理员", -1)
			return
		}
	}

	rds := cache.RedisPool.Get()
	defer rds.Close()

	// 判断是否开启在线客服
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminInfo.Id)
	var onlineKey, onlineToken string
	onlineKey = models.TokenParamsPrefix(models.HomePrefixTokenKey, settingAdminId)
	onlineUserInfo := models.NewUser(nil).AndWhere("admin_id=?", adminInfo.Id).AndWhere("type=?", models.UserTypeOnline).AndWhere("status=?", models.UserStatusActivate).FindOne()
	if onlineUserInfo != nil {
		onlineToken = router.TokenManager.Generate(rds, onlineKey, onlineUserInfo.AdminId, onlineUserInfo.Id)
	}

	tokenKey := models.TokenParamsPrefix(models.AdminPrefixTokenKey, adminInfo.Id)
	body.SuccessJSON(w, &loginData{
		Menu:       models.NewAdminMenu(nil).GetAdminMenuList(adminInfo.Id),
		RouterList: utils.GetMapValues(models.NewAdminAuthChild(nil).GetRolesRouteList(models.NewAdminAuthAssignment(nil).GetAdminRoleList(adminInfo.Id))),
		UserInfo: &userInfo{
			Id:         adminInfo.Id,
			Username:   adminInfo.UserName,
			Nickname:   adminInfo.Nickname,
			Email:      adminInfo.Email,
			Avatar:     adminInfo.Avatar,
			Money:      adminInfo.Money,
			Data:       adminInfo.Data,
			InviteCode: models.NewUserInvite(nil).GetInviteCode(adminInfo.Id, 0),
			ExpiredAt:  adminInfo.ExpiredAt,
			UpdatedAt:  adminInfo.UpdatedAt,
		},
		Token:       router.TokenManager.Generate(rds, tokenKey, adminInfo.Id, 0),
		TokenKey:    tokenKey,
		OnlineKey:   onlineKey,
		OnlineToken: onlineToken,
	})
}
