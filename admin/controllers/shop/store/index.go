package store

import (
	"basic/models"
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"net/http"
	"strings"
)

type indexParams struct {
	UserName            string             `json:"username"`             //用户名
	AdminName           string             `json:"admin_name"`           //管理员名
	SalesVolume         int64              `json:"sales_volume"`         //销售额
	VisitorCount        int64              `json:"visitor_count"`        //访客数
	OrderCount          int64              `json:"order_count"`          //订单数
	YesterdayDifference int64              `json:"yesterday_difference"` //昨日差
	Rating              int64              `json:"rating"`               //评分
	PendingPayment      string             `json:"pending_payment"`      //待付款
	PendingShipment     string             `json:"pending_shipment"`     //待发货
	PendingReceipt      string             `json:"pending_receipt"`      //待收货
	AfterSalesService   string             `json:"after_sales_service"`  //待售后
	PendingReview       string             `json:"pending_review"`       //待评论
	ShopName            string             `json:"shop_name"`            //店铺名称
	Phone               string             `json:"phone"`                //电话
	StoreType           string             `json:"store_type"`           //类型
	Keywords            string             `json:"keywords"`             //关键词
	Description         string             `json:"description"`          //描述
	Pagination          *define.Pagination `json:"pagination"`           //	分页
}

type indexData struct {
	models.StoreAttrs
	AdminName string `json:"admin_name"`
	UserName  string `json:"username"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取子级包括自己ID
	// 从token 中获取管理员ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	// 获取管理员ID，用于判断该管理员是否有权限获取信息
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)

	model := models.NewStore(nil)
	// 判断用户收藏中的管理员ID在管理员中
	model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	// 查询条件
	define.NewFilterEmpty(model.Db).
		// 添加product_name，date字段的条件
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
		String("shop_name=?", params.ShopName).
		String("phone=?", params.Phone).
		String("store_type=?", params.StoreType).
		String("keywords=?", params.Keywords).
		String("description=?", params.Description).
		Pagination(params.Pagination)

	// 管理员名称
	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}
	// 用户名称
	if params.UserName != "" {
		model.Db.AndWhere("user_id in (" + strings.Join(models.NewUser(nil).FindUserLikeNameIds(params.UserName), ",") + ")")
	}

	// 获取数据
	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.UserId, &tmp.AdminId, &tmp.SalesVolume, &tmp.VisitorCount, &tmp.OrderCount, &tmp.YesterdayDifference, &tmp.Rating, &tmp.PendingPayment, &tmp.PendingShipment, &tmp.PendingReceipt, &tmp.AfterSalesService, &tmp.PendingReview, &tmp.StoreLogo, &tmp.StoreName, &tmp.Phone, &tmp.StoreType, &tmp.Keywords, &tmp.Status, &tmp.Description)
		// 获取对应的管理员名
		adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}
		// 获取对应的用户名
		userInfo := models.NewUser(nil).AndWhere("id=?", tmp.UserId).FindOne()
		if userInfo != nil {
			tmp.UserName = userInfo.UserName
		}
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
