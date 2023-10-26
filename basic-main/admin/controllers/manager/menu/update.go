package menu

import (
	"basic/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type updateParams struct {
	Id     int64  `json:"id" validate:"required"`
	Parent int64  `json:"parent"`
	Name   string `json:"name"`
	Sort   int64  `json:"sort"`
	Route  string `json:"route"`
	Status int64  `json:"status" validate:"omitempty,oneof=-1 10"`
	Data   string `json:"data"`
}

// Update 菜单更新
func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updateParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	model := models.NewAdminMenu(nil)
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("name=?", params.Name).
		String("route=?", params.Route).
		String("data=?", params.Data).
		Int64("parent=?", params.Parent).
		Int64("sort=?", params.Sort).
		Int64("status=?", params.Status)

	_, err = model.AndWhere("id=?", params.Id).Update()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
