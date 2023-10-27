package index

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type updateParams struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// Update 当前管理员更新
func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updateParams)
	_ = body.ReadJSON(r, params)

	//  参数验证
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	//  获取请求管理员的uid
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	//  实例化模型
	model := models.NewAdminUser(nil)
	//  获取当前时间
	nowTime := time.Now()
	//  模型设置更新   过滤参数
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("email=?", params.Email).
		String("nickname=?", params.Nickname).
		String("avatar=?", params.Avatar).
		Int64("updated_at=?", nowTime.Unix())

	//  模型增加where条件并更新
	_, err = model.AndWhere("id = ?", adminId).Update()
	if err != nil {
		panic(err)
	}
	body.SuccessJSON(w, "ok")
}
