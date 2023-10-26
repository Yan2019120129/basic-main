package models

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// CommodityOrderAttrs 数据库模型属性
type CommodityOrderAttrs struct {
	Id                       int64   `json:"id"`                        //订单ID
	AdminId                  int64   `json:"admin_id"`                  //管理员ID
	ShopId                   int64   `json:"shop_id"`                   //店铺ID
	ShopLogo                 string  `json:"shop_logo"`                 //店铺LOG
	ShopName                 string  `json:"shop_name"`                 //店铺名
	ProductId                int64   `json:"product_id"`                //商品ID
	ProductImage             string  `json:"product_image"`             //商品图片
	ProductDescription       string  `json:"product_description"`       //商品描述
	AttributesSpecifications string  `json:"attributes_specifications"` //属性价格
	OriginalPrice            float64 `json:"original_price"`            //原价
	TransactionPrice         float64 `json:"transaction_price"`         //成交价
	Quantity                 int64   `json:"quantity"`                  //数量
	UserId                   int64   `json:"user_id"`                   //用户ID
	PaymentMethod            string  `json:"payment_method"`            //支付方式
	ShippingAddress          string  `json:"shipping_address"`          //收货地址
	UserOperate              int64   `json:"user_operate"`              //用户操作：1已下单，2取消订单
	OrderStatus              int64   `json:"order_status"`              //订单状态：0待处理，1已支付，2已发货，3已送达，-1已取消，-2已删除
	UserOperateTime          int64   `json:"user_operate_time"`         //用户操作时间
	StoreId                  int64   `json:"store_id"`                  //店家ID
	StoreOperate             int64   `json:"store_operate"`             //店家操作：1接受订单操作，2拒绝订单操作，3发货操作，4取消订单操作
	StoreOperateTime         int64   `json:"store_operate_time"`        //店家操作时间
	AdminOperateTime         int64   `json:"admin_operate_time"`        //管理员操作时间
}

const (
	UserOperatePlace      = 1  //	1已下单
	UserOperateCancel     = -1 //	2取消订单
	OrderStatusPending    = 0  //	0待处理
	OrderStatusPay        = 1  //	1已支付
	OrderStatusShipments  = 2  //	2已发货
	OrderStatusDelivery   = 3  //	3已送达
	OrderStatusCancel     = -1 //	-1已取消
	OrderStatusDelete     = -2 //	-2已删除
	StoreOperateReceive   = 1  //	1接受订单操作
	StoreOperateRefuse    = 2  //	2拒绝订单操作
	StoreOperateShipments = 3  //	3发货操作
	StoreOperateCancel    = 4  //	4取消订单操作
)

// CommodityOrder 数据库模型
type CommodityOrder struct {
	define.Db
}

// NewCommodityOrder 创建数据库模型
func NewCommodityOrder(tx *sql.Tx) *CommodityOrder {
	return &CommodityOrder{
		database.DbPool.NewDb(tx).Table("commodity_order"),
	}
}

// AndWhere where条件
func (c *CommodityOrder) AndWhere(str string, arg ...any) *CommodityOrder {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *CommodityOrder) FindOne() *CommodityOrderAttrs {
	attrs := new(CommodityOrderAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.ShopId, &attrs.ShopLogo, &attrs.ShopName, &attrs.ProductId, &attrs.ProductImage, &attrs.ProductDescription, &attrs.AttributesSpecifications, &attrs.OriginalPrice, &attrs.TransactionPrice, &attrs.Quantity, &attrs.UserId, &attrs.PaymentMethod, &attrs.ShippingAddress, &attrs.UserOperate, &attrs.OrderStatus, &attrs.UserOperateTime, &attrs.StoreId, &attrs.StoreOperate, &attrs.StoreOperateTime, &attrs.AdminOperateTime)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *CommodityOrder) FindMany() []*CommodityOrderAttrs {
	data := make([]*CommodityOrderAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(CommodityOrderAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.ShopId, &tmp.ShopLogo, &tmp.ShopName, &tmp.ProductId, &tmp.ProductImage, &tmp.ProductDescription, &tmp.AttributesSpecifications, &tmp.OriginalPrice, &tmp.TransactionPrice, &tmp.Quantity, &tmp.UserId, &tmp.PaymentMethod, &tmp.ShippingAddress, &tmp.UserOperate, &tmp.OrderStatus, &tmp.UserOperateTime, &tmp.StoreId, &tmp.StoreOperate, &tmp.StoreOperateTime, &tmp.AdminOperateTime)
		data = append(data, tmp)
	})
	return data
}
