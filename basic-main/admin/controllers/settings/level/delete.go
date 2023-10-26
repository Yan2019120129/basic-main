package level

import (
	"basic/models"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type deleteParams struct {
	Ids []int64 `json:"id" validate:"required"`
}

// Delete 等级删除
func Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(deleteParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	//  删除
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	for _, v := range params.Ids {
		model := models.NewUserLevel(nil)
		if adminId != models.AdminUserSupermanId {
			adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
			model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
		}
		_, err = model.Value("status=?").Args(models.UserLevelStatusDelete).
			AndWhere("id=?", v).Update()
		if err != nil {
			panic(err)
		}
	}
	body.SuccessJSON(w, "ok")
}
