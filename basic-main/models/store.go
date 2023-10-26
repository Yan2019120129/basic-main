package models

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// StoreAttrs 数据库模型属性
type StoreAttrs struct {
	Id                  int64  `json:"id"`                   //商店ID
	UserId              int64  `json:"user_id"`              //用户ID
	AdminId             int64  `json:"admin_id"`             //管理员ID
	SalesVolume         int64  `json:"sales_volume"`         //销售额
	VisitorCount        int64  `json:"visitor_count"`        //访客数
	OrderCount          int64  `json:"order_count"`          //订单数
	YesterdayDifference int64  `json:"yesterday_difference"` //昨日差
	Rating              int64  `json:"rating"`               //评分
	PendingPayment      string `json:"pending_payment"`      //待付款
	PendingShipment     string `json:"pending_shipment"`     //待发货
	PendingReceipt      string `json:"pending_receipt"`      //待收货
	AfterSalesService   string `json:"after_sales_service"`  //待售后
	PendingReview       string `json:"pending_review"`       //待评论
	StoreLogo           string `json:"store_logo"`           //店铺log
	StoreName           string `json:"store_name"`           //店铺名称
	Phone               string `json:"phone"`                //电话
	StoreType           string `json:"store_type"`           //类型
	Keywords            string `json:"keywords"`             //关键词
	Status              int64  `json:"status"`               //状态： -2关闭， -1整改维护， 1在使用， 10启用
	Description         string `json:"description"`          //描述
}

const (
	StoreStatusDelete   = -2 //-2关闭
	StoreStatusMaintain = -1 //-1整改维护
	StoreStatusUse      = 1  // 1在使用
	StoreStatusActive   = 10 // 10启用
)

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

// AndWhere where条件
func (c *Store) AndWhere(str string, arg ...any) *Store {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *Store) FindOne() *StoreAttrs {
	attrs := new(StoreAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.UserId, &attrs.AdminId, &attrs.SalesVolume, &attrs.VisitorCount, &attrs.OrderCount, &attrs.YesterdayDifference, &attrs.Rating, &attrs.PendingPayment, &attrs.PendingShipment, &attrs.PendingReceipt, &attrs.AfterSalesService, &attrs.PendingReview, &attrs.StoreLogo, &attrs.StoreName, &attrs.Phone, &attrs.StoreType, &attrs.Keywords, &attrs.Status, &attrs.Description)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *Store) FindMany() []*StoreAttrs {
	data := make([]*StoreAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(StoreAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.UserId, &tmp.AdminId, &tmp.SalesVolume, &tmp.VisitorCount, &tmp.OrderCount, &tmp.YesterdayDifference, &tmp.Rating, &tmp.PendingPayment, &tmp.PendingShipment, &tmp.PendingReceipt, &tmp.AfterSalesService, &tmp.PendingReview, &tmp.StoreLogo, &tmp.StoreName, &tmp.Phone, &tmp.StoreType, &tmp.Keywords, &tmp.Status, &tmp.Description)
		data = append(data, tmp)
	})
	return data
}
