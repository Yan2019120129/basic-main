package my_test

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// StoreAttrs 表示商店信息的结构体
type StoreAttrs struct {
	ID                  int64  `json:"id"`                   // 商店ID
	UserID              int64  `json:"user_id"`              // 用户ID
	SalesVolume         int    `json:"sales_volume"`         // 销售额
	VisitorCount        int    `json:"visitor_count"`        // 访客数
	OrderCount          int    `json:"order_count"`          // 订单数
	YesterdayDifference int    `json:"yesterday_difference"` // 昨日差
	Rating              int    `json:"rating"`               // 评分
	PendingPayment      string `json:"pending_payment"`      // 待付款
	PendingShipment     string `json:"pending_shipment"`     // 待发货
	PendingReceipt      string `json:"pending_receipt"`      // 待收货
	AfterSalesService   string `json:"after_sales_service"`  // 待售后
	PendingReview       string `json:"pending_review"`       // 待评论
	ShopLogo            string `json:"shop_logo"`            // 店铺log
	ShopName            string `json:"shop_name"`            // 店铺名称
	Phone               string `json:"phone"`                // 电话
	Type                string `json:"type"`                 // 类型
	Keywords            string `json:"keywords"`             // 关键词
	Description         string `json:"description"`          // 描述
}

// Store 数据库模型
type Store struct {
	define.Db
}

// NewStore 创建数据库模型
func NewStore(tx *sql.Tx) *Store {
	return &Store{
		database.DbPool.NewDb(tx).Table("store"),
	}
}

// AndWhere 添加 WHERE 子句
func (s *Store) AndWhere(str string, args ...interface{}) *Store {
	s.Db.AndWhere(str, args...)
	return s
}

// FindOne 查询单个商店信息
func (s *Store) FindOne() *StoreAttrs {
	store := new(StoreAttrs)
	s.QueryRow(func(row *sql.Row) {
		err := row.Scan(&store.ID, &store.UserID, &store.SalesVolume, &store.VisitorCount, &store.OrderCount, &store.YesterdayDifference,
			&store.Rating, &store.PendingPayment, &store.PendingShipment, &store.PendingReceipt, &store.AfterSalesService,
			&store.PendingReview, &store.ShopLogo, &store.ShopName, &store.Phone, &store.Type, &store.Keywords, &store.Description)
		if err != nil {
			store = nil
		}
	})
	return store
}

// FindMany 查询多个商店信息
func (s *Store) FindMany() []*StoreAttrs {
	data := make([]*StoreAttrs, 0)
	s.Query(func(rows *sql.Rows) {
		tmp := new(StoreAttrs)
		_ = rows.Scan(&tmp.ID, &tmp.UserID, &tmp.SalesVolume, &tmp.VisitorCount, &tmp.OrderCount, &tmp.YesterdayDifference,
			&tmp.Rating, &tmp.PendingPayment, &tmp.PendingShipment, &tmp.PendingReceipt, &tmp.AfterSalesService,
			&tmp.PendingReview, &tmp.ShopLogo, &tmp.ShopName, &tmp.Phone, &tmp.Type, &tmp.Keywords, &tmp.Description)
		data = append(data, tmp)
	})
	return data
}
