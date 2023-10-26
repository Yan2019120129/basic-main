package models

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// ShopStatisticalRecordAttrs 数据库模型属性
type ShopStatisticalRecordAttrs struct {
	Id                    int64   `json:"id"`                      //统计记录ID
	ShopId                int64   `json:"shop_id"`                 //店铺ID
	AdminId               int64   `json:"admin_id"`                //管理员ID
	VisitorCount          int64   `json:"visitor_count"`           //访客数
	OrderCount            int64   `json:"order_count"`             //订单数
	Earnings              float64 `json:"earnings"`                //收益
	ShopFavoritesCount    int64   `json:"shop_favorites_count"`    //店铺收藏量
	Credit                int64   `json:"credit"`                  //信用
	ProductFavoritesCount int64   `json:"product_favorites_count"` //商品收藏量
	ProductCount          int64   `json:"product_count"`           //商品数量
	Time                  int64   `json:"time"`                    //时间
}

// ShopStatisticalRecord 数据库模型
type ShopStatisticalRecord struct {
	define.Db
}

// NewShopStatisticalRecord 创建数据库模型
func NewShopStatisticalRecord(tx *sql.Tx) *ShopStatisticalRecord {
	return &ShopStatisticalRecord{
		database.DbPool.NewDb(tx).Table("shop_statistical_record"),
	}
}

// AndWhere where条件
func (c *ShopStatisticalRecord) AndWhere(str string, arg ...any) *ShopStatisticalRecord {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *ShopStatisticalRecord) FindOne() *ShopStatisticalRecordAttrs {
	attrs := new(ShopStatisticalRecordAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.ShopId, &attrs.VisitorCount, &attrs.OrderCount, &attrs.Earnings, &attrs.ShopFavoritesCount, &attrs.Credit, &attrs.ProductFavoritesCount, &attrs.ProductCount, &attrs.Time)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *ShopStatisticalRecord) FindMany() []*ShopStatisticalRecordAttrs {
	data := make([]*ShopStatisticalRecordAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(ShopStatisticalRecordAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.ShopId, &tmp.VisitorCount, &tmp.OrderCount, &tmp.Earnings, &tmp.ShopFavoritesCount, &tmp.Credit, &tmp.ProductFavoritesCount, &tmp.ProductCount, &tmp.Time)
		data = append(data, tmp)
	})
	return data
}
