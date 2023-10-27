package store

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/utils/body"
	"net/http"
)

type createParams struct {
	UserName            string `json:"user_name"  validate:"required"`             //用户名称
	AdminName           int64  `json:"admin_name"  validate:"required"`            //管理员名
	SalesVolume         int64  `json:"sales_volume"  validate:"required" `         //销售额
	VisitorCount        int64  `json:"visitor_count"  validate:"required" `        //访客数
	OrderCount          int64  `json:"order_count"  validate:"required"`           //订单数
	YesterdayDifference int64  `json:"yesterday_difference"  validate:"required" ` //昨日差
	Rating              int64  `json:"rating"  validate:"required" `               //评分
	PendingPayment      string `json:"pending_payment"  validate:"required" `      //待付款
	PendingShipment     string `json:"pending_shipment"  validate:"required" `     //待发货
	PendingReceipt      string `json:"pending_receipt"  validate:"required" `      //待收货
	AfterSalesService   string `json:"after_sales_service"  validate:"required" `  //待售后
	PendingReview       string `json:"pending_review"  validate:"required" `       //待评论
	StoreLogo           string `json:"store_logo"  validate:"required" `           //店铺log
	StoreName           string `json:"store_name"  validate:"required" `           //店铺名称
	Phone               string `json:"phone"  validate:"required" `                //电话
	StoreType           string `json:"store_type"  validate:"required" `           //类型
	Keywords            string `json:"keywords"  validate:"required" `             //关键词
	Status              int64  `json:"status"  validate:"required" `               //状态： -2关闭， -1整改维护， 1在使用， 10启用
	Description         string `json:"description"  validate:"required" `          //描述
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)
	tx := database.DbPool.GetTx()
	defer tx.Rollback()
	model := models.NewStore(nil)
	model.Field("image", "name", "status").
		Args(params.Image, params.Name, params.Status).
		Insert()
	//	TODO... 消费账单...
	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
