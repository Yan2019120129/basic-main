package order

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type createParams struct {
	UserName  string `json:"username" validate:"required"`
	ProductId int64  `json:"product_id" validate:"required,gt=0"`
	Nums      int64  `json:"nums" validate:"required,gt=0"`
	Type      int64  `json:"type" validate:"required,oneof=1"`
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	//  验证参数
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 查询用户说会否存在
	userInfo := models.NewUser(nil).AndWhere("username=?", params.UserName).FindOne()
	if userInfo == nil {
		body.ErrorJSON(w, "用户不存在", -1)
		return
	}

	// 查询产品Id是否存在
	productInfo := models.NewProduct(nil).AndWhere("id=?", params.ProductId).AndWhere("status>?", models.ProductStatusDelete).FindOne()
	if productInfo == nil {
		body.ErrorJSON(w, "产品不存在", -1)
		return
	}

	// 如果是超级管理员能修改所有用户， 不是超级管理员只能修改自己用户
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	nowTime := time.Now()

	if adminId != models.AdminUserSupermanId && adminId != userInfo.AdminId {
		body.ErrorJSON(w, "权限不足", -1)
		return
	}

	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	_, err = models.NewProductOrder(tx).
		Field("admin_id", "user_id", "product_id", "order_sn", "money", "type", "nums", "data", "expired_at", "updated_at", "created_at").
		Args(userInfo.AdminId, userInfo.Id, params.ProductId, utils.NewRandom().OrderSn(), productInfo.Money, params.Type, params.Nums, productInfo.Data, nowTime.Unix(), nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	//	TODO... 消费账单...
	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
