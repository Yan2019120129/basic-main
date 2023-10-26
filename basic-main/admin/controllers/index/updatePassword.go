package index

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type updatePasswordParams struct {
	Type        int64  `json:"type" validate:"required,oneof=1 2"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

const (
	// UpdateLoginPassword 更新登陆密码
	UpdateLoginPassword = 1
	// UpdateSecurityPassword 更新安全密钥
	UpdateSecurityPassword = 2
)

// UpdatePassword 当前管理员更新密码
func UpdatePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updatePasswordParams)
	_ = body.ReadJSON(r, params)

	//  参数验证
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	//  获取管理员ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	adminInfo := models.NewAdminUser(nil).AndWhere("id=?", adminId).FindOne()

	//  实例化模型
	model := models.NewAdminUser(nil)
	//  获取当前时间
	nowTime := time.Now()

	if params.Type == UpdateLoginPassword {
		if adminInfo == nil || adminInfo.Password != utils.PasswordEncrypt(params.OldPassword) {
			body.ErrorJSON(w, "旧密码不正确", -1)
			return
		}
		define.NewFilterEmpty(model.Db).SetUpdateOpt().
			String("password=?", utils.PasswordEncrypt(params.NewPassword)).
			Int64("updated_at=?", nowTime.Unix())
	}

	if params.Type == UpdateSecurityPassword {
		if adminInfo == nil || adminInfo.SecurityKey != utils.PasswordEncrypt(params.OldPassword) {
			body.ErrorJSON(w, "旧密码不正确", -1)
			return
		}

		define.NewFilterEmpty(model.Db).SetUpdateOpt().
			String("security_key=?", utils.PasswordEncrypt(params.NewPassword)).
			Int64("updated_at=?", nowTime.Unix())
	}
	//  模型增加where条件并更新
	_, err = model.AndWhere("id=?", adminId).Update()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
