package my_test

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// StatisticsRecordAttrs 表示统计记录信息的结构体
type StatisticsRecordAttrs struct {
	ID                    int64   `json:"id"`                                // 统计记录ID
	ShopID                int64   `json:"shop_id"`                           // 店铺ID
	VisitorCount          int     `json:"visitor_count"`                     // 访客数
	OrderCount            int     `json:"order_count"`                       // 订单数
	Earnings              float64 `json:"earnings" sql:"type:decimal(12,2)"` // 收益
	ShopFavoritesCount    int     `json:"shop_favorites_count"`              // 店铺收藏量
	Credit                int     `json:"credit"`                            // 信用
	ProductFavoritesCount int     `json:"product_favorites_count"`           // 商品收藏量
	ProductCount          int     `json:"product_count"`                     // 商品数量
	Time                  int     `json:"time"`                              // 时间
}

// StatisticsRecord 数据库模型
type StatisticsRecord struct {
	define.Db
}

// NewStatisticsRecord 创建数据库模型
func NewStatisticsRecord(tx *sql.Tx) *StatisticsRecord {
	return &StatisticsRecord{
		database.DbPool.NewDb(tx).Table("statistics_record"),
	}
}

// AndWhere 添加 WHERE 子句
func (s *StatisticsRecord) AndWhere(str string, args ...interface{}) *StatisticsRecord {
	s.Db.AndWhere(str, args...)
	return s
}

// FindOne 查询单个统计记录信息
func (s *StatisticsRecord) FindOne() *StatisticsRecordAttrs {
	record := new(StatisticsRecordAttrs)
	s.QueryRow(func(row *sql.Row) {
		err := row.Scan(&record.ID, &record.ShopID, &record.VisitorCount, &record.OrderCount, &record.Earnings, &record.ShopFavoritesCount,
			&record.Credit, &record.ProductFavoritesCount, &record.ProductCount, &record.Time)
		if err != nil {
			record = nil
		}
	})
	return record
}

// FindMany 查询多个统计记录信息
func (s *StatisticsRecord) FindMany() []*StatisticsRecordAttrs {
	data := make([]*StatisticsRecordAttrs, 0)
	s.Query(func(rows *sql.Rows) {
		tmp := new(StatisticsRecordAttrs)
		_ = rows.Scan(&tmp.ID, &tmp.ShopID, &tmp.VisitorCount, &tmp.OrderCount, &tmp.Earnings, &tmp.ShopFavoritesCount,
			&tmp.Credit, &tmp.ProductFavoritesCount, &tmp.ProductCount, &tmp.Time)
		data = append(data, tmp)
	})
	return data
}
