package attention

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
	"net/http"
	"strings"
	"time"
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
	nowTime := time.Now()
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	for _, v := range params.Ids {
		// 判断是否有修改权限
		model := models.NewUserAttention(nil)
		if adminId != models.AdminUserSupermanId {
			adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
			model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
		}
		// 修改用户关注店铺信息
		_, err = model.Value("status=?", "date=?").Args(models.UserAttentionDelete, nowTime.Unix()).
			AndWhere("id=?", v).Update()
		if err != nil {
			panic(err)
		}
	}
	body.SuccessJSON(w, "ok")
}
