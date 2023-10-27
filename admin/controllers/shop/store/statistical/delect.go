package statistical

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
	"net/http"
	"strings"
)

type deleteParams struct {
	Ids []int64 `json:"id" validate:"required"`
}

// Delete 删除产品
func Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(deleteParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	adminId := router.TokenManager.GetContextClaims(r).AdminId
	tx := database.DbPool.GetTx()
	defer tx.Rollback()
	for _, v := range params.Ids {
		// 判断是否有修改权限
		model := models.NewShopStatisticalRecord(nil)
		if adminId != models.AdminUserSupermanId {
			adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
			model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
		}
		// 修改用户地址状态为删除信息
		_, err = model.AndWhere("id=?", v).Delete()
		if err != nil {
			body.ErrorJSON(w, err.Error(), -1)
			return
		}
	}
	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
