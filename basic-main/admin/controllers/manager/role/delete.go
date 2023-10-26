package role

import (
	"basic/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type deleteParams struct {
	NameList []string `json:"name" validate:"required"`
}

// Delete 删除角色
func Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(deleteParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	for _, name := range params.NameList {
		// 删除角色
		_, err = models.NewAdminAuthItem(tx).
			AndWhere("name=?", name).AndWhere("type=?", models.AdminAuthItemTypeManage).Delete()
		if err != nil {
			panic(err)
		}

		_, err = models.NewAdminAuthChild(tx).AndWhere("parent=?", name).AndWhere("type=?", models.AdminAuthItemTypeRouteName).Delete()
		if err != nil {
			panic(err)
		}
	}

	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
