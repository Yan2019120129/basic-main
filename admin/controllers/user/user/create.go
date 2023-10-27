package user

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type createParams struct {
	UserName    string `json:"username" validate:"required,min=4"`
	Password    string `json:"password" validate:"required,min=6"`
	SecurityKey string `json:"security_key" validate:"required,min=6"`
	Type        int64  `json:"type" validate:"required,oneof=-10 -1 10"`
}

// Create 新增用户
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 如果该用户存在,那么提示
	oldUserInfo := models.NewUser(nil).AndWhere("username=?", params.UserName).FindOne()
	if oldUserInfo != nil {
		body.ErrorJSON(w, "用户名已存在", -1)
		return
	}

	//  router
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	nowTime := time.Now()
	_, err = models.NewUser(nil).
		Field("admin_id", "username", "password", "security_key", "type", "created_at", "updated_at").
		Args(adminId, params.UserName, utils.PasswordEncrypt(params.Password), utils.PasswordEncrypt(params.SecurityKey), params.Type, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
