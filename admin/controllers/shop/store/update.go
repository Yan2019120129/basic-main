package store

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
	Id                  int64  `json:"id" validate:"required,gt=0"` //类目ID
	SalesVolume         int64  `json:"sales_volume"`                //销售额
	VisitorCount        int64  `json:"visitor_count"`               //访客数
	OrderCount          int64  `json:"order_count"`                 //订单数
	YesterdayDifference int64  `json:"yesterday_difference"`        //昨日差
	Rating              int64  `json:"rating"`                      //评分
	PendingPayment      string `json:"pending_payment"`             //待付款
	PendingShipment     string `json:"pending_shipment"`            //待发货
	PendingReceipt      string `json:"pending_receipt"`             //待收货
	AfterSalesService   string `json:"after_sales_service"`         //待售后
	PendingReview       string `json:"pending_review"`              //待评论
	ShopLogo            string `json:"shop_logo"`                   //店铺log
	ShopName            string `json:"shop_name"`                   //店铺名称
	Phone               string `json:"phone"`                       //电话
	StoreType           string `json:"store_type"`                  //类型
	Keywords            string `json:"keywords"`                    //关键词
	Status              int64  `json:"status"`                      //状态： -2关闭， -1整改维护， 1在使用， 10启用
	Description         string `json:"description"`                 //描述
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
	model := models.NewStore(nil)

	// 开启事务
	tx := database.DbPool.GetTx()
	defer tx.Rollback()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		Int64("sales_volume=?", params.SalesVolume).
		Int64("visitor_count=?", params.VisitorCount).
		Int64("order_count=?", params.OrderCount).
		Int64("yesterday_difference=?", params.YesterdayDifference).
		Int64("rating=?", params.Rating).
		String("pending_payment=?", params.PendingPayment).
		String("pending_shipment=?", params.PendingShipment).
		String("pending_receipt=?", params.PendingReceipt).
		String("after_sales_service=?", params.AfterSalesService).
		String("pending_review=?", params.PendingReview).
		String("shop_logo=?", params.ShopLogo).
		String("shop_name=?", params.ShopName).
		String("phone=?", params.Phone).
		String("store_type=?", params.StoreType).
		String("Keywords=?", params.Keywords).
		Int64("Status=?", params.Status).
		String("Description=?", params.Description)
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
