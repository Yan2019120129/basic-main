package level

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type createParams struct {
	Name  string  `json:"name" validate:"required"`
	Icon  string  `json:"icon" validate:"required"`
	Money float64 `json:"money"`
	Days  int64   `json:"days"`
	Data  string  `json:"data"`
}

// Create 新增管理员
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	//  验证参数
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	adminId := router.TokenManager.GetContextClaims(r).AdminId
	nowTime := time.Now()
	var userLevel int64 = 1
	var levelDays int64 = -1
	if params.Days > 0 {
		levelDays = params.Days
	}
	endLevelModel := models.NewUserLevel(nil)
	endLevelModel.AndWhere("admin_id=?", adminId).AndWhere("status=?", models.UserLevelStatusActivate).OrderBy("level desc")
	endLevelInfo := endLevelModel.FindOne()
	if endLevelInfo != nil {
		userLevel = endLevelInfo.Level + 1
	}

	//  模型插入数据
	_, err = models.NewUserLevel(nil).
		Field("admin_id", "name", "icon", "level", "money", "days", "created_at").
		Args(adminId, params.Name, params.Icon, userLevel, params.Money, levelDays, nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
