package lang

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
	Name   string `json:"name" validate:"max=50"`
	Alias  string `json:"alias" validate:"max=50"`
	Icon   string `json:"icon" validate:"max=255"`
	Status int64  `json:"status" validate:"omitempty,oneof=-1 10"`
	Sort   int64  `json:"sort"`
	Data   string `json:"data" validate:"max=255"`
}

// Update 更新国家语言
func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updateParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	//  实例化模型
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	model := models.NewLang(nil)
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("name=?", params.Name).
		String("alias=?", params.Alias).
		String("icon=?", params.Icon).
		String("data=?", params.Data).
		Int64("sort=?", params.Sort).
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
