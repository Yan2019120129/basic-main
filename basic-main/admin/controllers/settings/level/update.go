package level

import (
	"basic/models"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type updateParams struct {
	Id     int64   `json:"id" validate:"required"`
	Name   string  `json:"name"`
	Icon   string  `json:"icon"`
	Level  int64   `json:"level"`
	Money  float64 `json:"money"`
	Days   int64   `json:"days"`
	Status int64   `json:"status"`
	Data   string  `json:"data"`
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

	adminId := router.TokenManager.GetContextClaims(r).AdminId
	nowTime := time.Now()
	//  实例化模型
	model := models.NewUserLevel(nil)
	var userLevel int64 = 0
	if params.Level > 0 {
		levelInfo := models.NewUserLevel(nil).AndWhere("admin_id=?", adminId).AndWhere("level=?", params.Level).FindOne()
		if levelInfo == nil {
			userLevel = params.Level
		}
	}

	//  模型设置更新   过滤参数
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("name=?", params.Name).
		String("icon=?", params.Icon).
		Int64("level=?", userLevel).
		Float64("money=?", params.Money).
		String("data=?", params.Data).
		Int64("days=?", params.Days).
		Int64("status=?", params.Status).
		Int64("updated_at=?", nowTime.Unix())

	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id = ?", params.Id).Update()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
