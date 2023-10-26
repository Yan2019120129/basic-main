package assets

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
	UserName string `json:"username" validate:"required"`
	AssetsId int64  `json:"assets_id" validate:"required"`
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	adminId := router.TokenManager.GetContextClaims(r).AdminId
	userInfo := models.NewUser(nil).AndWhere("username=?", params.UserName).AndWhere("admin_id=?", adminId).FindOne()
	if userInfo == nil {
		body.ErrorJSON(w, "用户名不存在", -1)
		return
	}

	// 如果资产已存在, 那么不能添加
	userAssetsInfo := models.NewUserAssets(nil).AndWhere("user_id=?", userInfo.Id).AndWhere("assets_id=?", params.AssetsId).AndWhere("status>?", models.ProductStatusDelete).FindOne()
	if userAssetsInfo != nil {
		body.ErrorJSON(w, "当前用户资产已存在", -1)
		return
	}

	nowTime := time.Now()
	_, err = models.NewUserAssets(nil).Field("admin_id", "user_id", "assets_id", "created_at", "updated_at").
		Args(userInfo.AdminId, userInfo.Id, params.AssetsId, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
