package shop_order

import (
	"basic/models"
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"net/http"
	"strings"
	"time"
)

type indexParams struct {
	UserName                 string                 `json:"username"`                  //用户名
	AdminName                string                 `json:"admin_name"`                //管理员名
	ShopName                 string                 `json:"shop_name"`                 //店铺名
	ProductDescription       string                 `json:"product_description"`       //商品描述
	AttributesSpecifications string                 `json:"attributes_specifications"` //属性价格
	OriginalPrice            float64                `json:"original_price"`            //原价
	TransactionPrice         float64                `json:"transaction_price"`         //成交价
	Quantity                 int64                  `json:"quantity"`                  //数量
	PaymentMethod            string                 `json:"payment_method"`            //支付方式
	ShippingAddress          string                 `json:"shipping_address"`          //收货地址
	UserOperate              int64                  `json:"user_operate"`              //用户操作：1已下单，2取消订单
	OrderStatus              int64                  `json:"order_status"`              //订单状态：0待处理，1已支付，2已发货，3已送达，-1已取消
	StoreOperate             int64                  `json:"store_operate"`             //店家操作：1接受订单操作，2拒绝订单操作，3发货操作，4取消订单操作
	UserOperateTime          *define.RangeTimeParam `json:"user_operate_time"`         // 时间戳
	StoreOperateTime         *define.RangeTimeParam `json:"store_operate_time"`        // 时间戳
	Pagination               *define.Pagination     `json:"pagination"`                //	分页
}

type indexData struct {
	models.CommodityOrderAttrs
	AdminName   string `json:"admin_name"`
	UserName    string `json:"username"`
	ProductName int64  `json:"product_name"` //商品名
	StoreName   int64  `json:"store_name"`   //店家名

}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取子级包括自己ID
	// 从token 中获取管理员ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	// 获取管理员ID，用于判断该管理员是否有权限获取信息
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(adminId, "site_timezone"))

	model := models.NewCommodityOrder(nil)
	// 判断用户收藏中的管理员ID在管理员中
	model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	// 查询条件
	define.NewFilterEmpty(model.Db).
		// 添加product_name，date字段的条件
		String("shop_name=?", params.ShopName).
		String("product_description=?", params.ProductDescription).
		String("attributes_specifications=?", params.AttributesSpecifications).
		Float64("original_price=?", params.OriginalPrice).
		Float64("transaction_price=?", params.TransactionPrice).
		Int64("quantity=?", params.Quantity).
		String("payment_method=?", params.PaymentMethod).
		String("shipping_address=?", params.ShippingAddress).
		Int64("user_operate=?", params.UserOperate).
		Int64("order_status=?", params.OrderStatus).
		Int64("store_operate=?", params.StoreOperate).
		RangeTime("user_operate_time between ? and ?", params.UserOperateTime, location).
		RangeTime("store_operate_time between ? and ?", params.StoreOperateTime, location).
		Pagination(params.Pagination)

	// 管理员名称ƒ
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
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.ShopId, &tmp.ShopLogo, &tmp.ShopName, &tmp.ProductId, &tmp.ProductImage, &tmp.ProductDescription, &tmp.AttributesSpecifications, &tmp.OriginalPrice, &tmp.TransactionPrice, &tmp.Quantity, &tmp.UserId, &tmp.PaymentMethod, &tmp.ShippingAddress, &tmp.UserOperate, &tmp.OrderStatus, &tmp.UserOperateTime, &tmp.StoreId, &tmp.StoreOperate, &tmp.StoreOperateTime)
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
