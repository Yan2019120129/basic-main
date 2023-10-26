package my_test

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// OrderAttrs 表示订单信息的结构体
type OrderAttrs struct {
	ID                 int64   `json:"id"`                                         // 订单ID
	ShopID             int64   `json:"shop_id"`                                    // 店铺ID
	ShopLogo           string  `json:"shop_logo"`                                  // 店铺LOG
	ShopName           string  `json:"shop_name"`                                  // 店铺名
	ProductID          int64   `json:"product_id"`                                 // 商品ID
	ProductImage       string  `json:"product_image"`                              // 商品图片
	ProductDescription string  `json:"product_description"`                        // 商品描述
	AttributesSpecs    string  `json:"attributes_specifications"`                  // 属性规格
	OriginalPrice      float64 `json:"original_price" sql:"type:decimal(10,2)"`    // 原价
	TransactionPrice   float64 `json:"transaction_price" sql:"type:decimal(10,2)"` // 成交价
	Quantity           int64   `json:"quantity"`                                   // 数量
	UserID             int64   `json:"user_id"`                                    // 用户ID
	PaymentMethod      string  `json:"payment_method"`                             // 支付方式
	ShoppingAddress    string  `json:"shopping_address"`                           // 收货地址
	UserOperate        int     `json:"user_operate"`                               // 用户操作
	OrderStatus        int     `json:"order_status"`                               // 订单状态
	UserOperateTime    int     `json:"user_operate_time"`                          // 用户操作时间
	StoreID            int64   `json:"store_id"`                                   // 店家ID
	StoreOperate       int     `json:"store_operate"`                              // 店家操作
	StoreOperateTime   int     `json:"store_operate_time"`                         // 店家操作时间
}

// OrderStatus 表示订单状态
const (
	OrderStatusPending   = 0  // 待处理
	OrderStatusPaid      = 1  // 已支付
	OrderStatusShipped   = 2  // 已发货
	OrderStatusDelivered = 3  // 已送达
	OrderStatusCanceled  = -1 // 已取消
)

// UserOperate 表示用户操作
const (
	UserOperatePlaceOrder  = 1 // 下单操作
	UserOperateCancelOrder = 2 // 取消订单操作
)

// StoreOperate 表示店家操作
const (
	StoreOperateAcceptOrder = 1 // 接受订单操作
	StoreOperateRejectOrder = 2 // 拒绝订单操作
	StoreOperateShipOrder   = 3 // 发货操作
	StoreOperateCancelOrder = 4 // 取消订单操作
)

// Order 数据库模型
type Order struct {
	define.Db
}

// NewOrder 创建数据库模型
func NewOrder(tx *sql.Tx) *Order {
	return &Order{
		database.DbPool.NewDb(tx).Table("order"),
	}
}

// AndWhere 添加 WHERE 子句
func (o *Order) AndWhere(str string, args ...interface{}) *Order {
	o.Db.AndWhere(str, args...)
	return o
}

// FindOne 查询单个订单信息
func (o *Order) FindOne() *OrderAttrs {
	order := new(OrderAttrs)
	o.QueryRow(func(row *sql.Row) {
		err := row.Scan(&order.ID, &order.ShopID, &order.ShopLogo, &order.ShopName, &order.ProductID, &order.ProductImage, &order.ProductDescription,
			&order.AttributesSpecs, &order.OriginalPrice, &order.TransactionPrice, &order.Quantity, &order.UserID, &order.PaymentMethod, &order.ShoppingAddress,
			&order.UserOperate, &order.OrderStatus, &order.UserOperateTime, &order.StoreID, &order.StoreOperate, &order.StoreOperateTime)
		if err != nil {
			order = nil
		}
	})
	return order
}

// FindMany 查询多个订单信息
func (o *Order) FindMany() []*OrderAttrs {
	data := make([]*OrderAttrs, 0)
	o.Query(func(rows *sql.Rows) {
		tmp := new(OrderAttrs)
		_ = rows.Scan(&tmp.ID, &tmp.ShopID, &tmp.ShopLogo, &tmp.ShopName, &tmp.ProductID, &tmp.ProductImage, &tmp.ProductDescription,
			&tmp.AttributesSpecs, &tmp.OriginalPrice, &tmp.TransactionPrice, &tmp.Quantity, &tmp.UserID, &tmp.PaymentMethod, &tmp.ShoppingAddress,
			&tmp.UserOperate, &tmp.OrderStatus, &tmp.UserOperateTime, &tmp.StoreID, &tmp.StoreOperate, &tmp.StoreOperateTime)
		data = append(data, tmp)
	})
	return data
}
