package shop_order

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
	"net/http"
	"time"
)

type createParams struct {
	AdminName                string  `json:"admin_name" validate:"required"`                //管理员名称
	ShopLogo                 string  `json:"shop_logo" validate:"required"`                 //店铺LOG
	ShopId                   int64   `json:"shop_id" validate:"required"`                   //店铺ID
	ShopName                 string  `json:"shop_name" validate:"required"`                 //店铺名
	ProductId                int64   `json:"product_id" validate:"required"`                //商品ID
	ProductImage             string  `json:"product_image" validate:"required"`             //商品图片
	ProductDescription       string  `json:"product_description" validate:"required"`       //商品描述
	AttributesSpecifications string  `json:"attributes_specifications" validate:"required"` //属性价格
	OriginalPrice            float64 `json:"original_price" validate:"required"`            //原价
	TransactionPrice         float64 `json:"transaction_price" validate:"required"`         //成交价
	Quantity                 int64   `json:"quantity" validate:"required"`                  //数量
	UserName                 string  `json:"user_name" validate:"required"`                 //用户名
	PaymentMethod            string  `json:"payment_method" validate:"required"`            //支付方式
	ShippingAddress          string  `json:"shipping_address" validate:"required"`          //收货地址
	OrderStatus              int64   `json:"order_status" validate:"required"`              //订单状态：0待处理，1已支付，2已发货，3已送达，-1已取消，-2已删除
	StoreId                  int64   `json:"store_id"`                                      //店家ID
	StoreName                string  `json:"store_name" validate:"required"`                //店家名
	StoreOperate             int64   `json:"store_operate" validate:"required,gt=0"`        //店家操作：1接受订单操作，2拒绝订单操作，3发货操作，4取消订单操作
	AdminOperateTime         int64   `json:"admin_operate_time" validate:"required"`        //管理员操作时间
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
	// 查询用户是否存在
	userInfo := models.NewUser(nil).AndWhere("username=?", params.UserName).FindOne()
	if userInfo == nil {
		body.ErrorJSON(w, "用户不存在", -1)
		return
	}
	// 查询产品Id是否存在
	productInfo := models.NewCommodity(nil).AndWhere("id=?", params.ProductId).AndWhere("status=?", models.CommodityStatusOnSale).FindOne()
	if productInfo == nil {
		body.ErrorJSON(w, "产品不存在", -1)
		return
	}
	// 查询商店Id是否存在
	storeInfo := models.NewStore(nil).AndWhere("id=?", params.ShopId).AndWhere("status>?", models.StoreStatusMaintain).FindOne()
	if storeInfo == nil {
		body.ErrorJSON(w, "商店不存在", -1)
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
	_, err = models.NewCommodityOrder(nil).
		Field("admin_id", "user_id", "shop_logo", "shop_id", "product_id", "order_sn", "product_image", "product_description", "attributes_specifications", "original_price", "transaction_price", "quantity", "payment_method", "shipping_address", "order_status", "admin_operate_time", "store_id", "shop_name", "user_operate_time").
		Args(userInfo.AdminId, userInfo.Id, params.ShopLogo, params.ShopId, params.ProductId, utils.NewRandom().OrderSn(), params.ProductImage, params.ProductDescription, params.AttributesSpecifications, params.OriginalPrice, params.TransactionPrice, params.Quantity, params.PaymentMethod, params.ShippingAddress, params.OrderStatus, nowTime.Unix(), params.StoreId, params.ShopName, nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}
	//	TODO... 消费账单...
	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
