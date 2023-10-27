package favorites

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
	"net/http"
	"strings"
)

type updateParams struct {
	Id           int64  `json:"id" validate:"required"`            //收藏ID
	UserId       int64  `json:"user_id" validate:"required"`       //用户ID
	AdminId      int64  `json:"admin_id"`                          //管理员ID
	CommodityId  int64  `json:"commodity_id" 	validate:"required"` //商品ID
	ProductName  string `json:"product_name" validate:"required"`  //商品名称
	ProductImage string `json:"product_imag" validate:"required"`  //商品图片
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
	// 从token获取管理员ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	model := models.NewUserFavorites(nil)

	// 开启事务
	tx := database.DbPool.GetTx()
	defer tx.Rollback()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("product_imag=?", params.ProductName).
		String("product_name=?", params.ProductImage)
	// 是否是超级管理员，不是没有权限修改
	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id = ?", params.Id).Update()
	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
