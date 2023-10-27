package chat

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type registerParams struct {
	Domain string `json:"domain" validate:"required"` //	目标域名
	UUID   string `json:"uuid" validate:"required"`   //	浏览器唯一标识
}

type registerData struct {
	Token    string `json:"token"`
	TokenKey string `json:"tokenKey"`
}

// Register 注册临时用户
func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(registerParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	rds := cache.RedisPool.Get()
	defer rds.Close()

	adminId := models.NewAdminUser(nil).FindDomainAdminId(params.Domain)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	var onlineKey, onlineToken string
	onlineKey = models.TokenParamsPrefix(models.HomePrefixTokenKey, settingAdminId)

	temporaryInfo := models.NewUser(nil).AndWhere("username=?", params.UUID).AndWhere("type=?", models.UserTypeTemporary).FindOne()
	if temporaryInfo != nil {
		onlineToken = router.TokenManager.Generate(rds, onlineKey, temporaryInfo.AdminId, temporaryInfo.Id)
	}

	//	创建临时用户
	if temporaryInfo == nil {
		nowTime := time.Now()
		userIP := utils.GetUserRealIP(r)
		currentUserId, err := models.NewUser(nil).Field("admin_id", "username", "type", "status", "ip4", "created_at", "updated_at").
			Value("?", "?", "?", "?", "INET_ATON(?)", "?", "?").
			Args(adminId, params.UUID, models.UserTypeTemporary, models.UserStatusActivate, userIP, nowTime.Unix(), nowTime.Unix()).
			Insert()
		if err != nil {
			panic(err)
		}
		onlineToken = router.TokenManager.Generate(rds, onlineKey, adminId, currentUserId)
	}

	body.SuccessJSON(w, &registerData{
		Token:    onlineToken,
		TokenKey: onlineKey,
	})
}
