package statistical

import (
	"basic/models"
	"fmt"
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
	Id                    int64   `json:"id" validate:"required,gt=0"` //类目ID
	VisitorCount          int64   `json:"visitor_count"`               //访客数
	OrderCount            int64   `json:"order_count"`                 //订单数
	Earnings              float64 `json:"earnings"`                    //收益
	ShopFavoritesCount    int64   `json:"shop_favorites_count"`        //店铺收藏量
	Credit                int64   `json:"credit"`                      //信用
	ProductFavoritesCount int64   `json:"product_favorites_count"`     //商品收藏量
	ProductCount          int64   `json:"product_count"`               //商品数量
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
	model := models.NewShopStatisticalRecord(nil)

	// 开启事务
	tx := database.DbPool.GetTx()
	defer tx.Rollback()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		Int64("visitor_count=?", params.VisitorCount).
		Int64("order_count=?", params.OrderCount).
		Float64("earnings=?", params.Earnings).
		Int64("shop_favorites_count=?", params.ShopFavoritesCount).
		Int64("credit=?", params.Credit).
		Int64("product_favorites_count=?", params.ProductFavoritesCount).
		Int64("product_count=?", params.ProductCount)
	// 是否是超级管理员，不是没有权限修改
	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id = ?", params.Id).Update()
	fmt.Println("err", err)
	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
