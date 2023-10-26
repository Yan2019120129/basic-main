package assets

import (
	"basic/models"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type updateParams struct {
	Id     int64  `json:"id" validate:"required"`
	Name   string `json:"name" validate:"omitempty,max=50"`
	Icon   string `json:"icon"`
	Type   int64  `json:"type" validate:"omitempty,oneof=1 2 3"`
	Status int64  `json:"status" validate:"omitempty,oneof=-1 10"`
	Data   string `json:"data"`
}

// Update 更新资产
func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updateParams)
	_ = body.ReadJSON(r, params)
	//  参数验证
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	//  实例化模型
	model := models.NewAssets(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("name=?", params.Name).
		String("icon=?", params.Icon).
		String("data=?", params.Data).
		Int64("type=?", params.Type).
		Int64("status=?", params.Status)

	//  模型增加where条件并更新
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
