package country

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
	LangId int64  `json:"lang_id"`
	Name   string `json:"name" validate:"omitempty,max=50"`
	Alias  string `json:"alias" validate:"omitempty,max=50"`
	Iso1   string `json:"iso1" validate:"omitempty,max=50"`
	Icon   string `json:"icon" validate:"omitempty,max=255"`
	Code   string `json:"code" validate:"omitempty,max=50"`
	Sort   int64  `json:"sort"`
	Status int64  `json:"status" validate:"omitempty,oneof=-1 10"`
	Data   string `json:"data" validate:"omitempty,max=255"`
}

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
	model := models.NewCountry(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		Int64("lang_id=?", params.LangId).
		String("name=?", params.Name).
		String("alias=?", params.Alias).
		String("iso1=?", params.Iso1).
		String("icon=?", params.Icon).
		String("code=?", params.Code).
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
