package my_test

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// FinancialStatisticsAttrs 表示财务统计信息的结构体
type FinancialStatisticsAttrs struct {
	ID             int64   `json:"id"`              // 主键
	ProductID      int64   `json:"product_id"`      // 商品ID
	Profit         float64 `json:"profit"`          // 利润
	UnitPrice      float64 `json:"unit_price"`      // 单价
	WholesalePrice float64 `json:"wholesale_price"` // 批发价
}

// FinancialStatistics 数据库模型
type FinancialStatistics struct {
	define.Db
}

// NewFinancialStatistics 创建数据库模型
func NewFinancialStatistics(tx *sql.Tx) *FinancialStatistics {
	return &FinancialStatistics{
		database.DbPool.NewDb(tx).Table("financial_statistics"),
	}
}

// AndWhere 添加 WHERE 子句
func (f *FinancialStatistics) AndWhere(str string, args ...interface{}) *FinancialStatistics {
	f.Db.AndWhere(str, args...)
	return f
}

// FindOne 查询单个财务统计信息
func (f *FinancialStatistics) FindOne() *FinancialStatisticsAttrs {
	stats := new(FinancialStatisticsAttrs)
	f.QueryRow(func(row *sql.Row) {
		err := row.Scan(&stats.ID, &stats.ProductID, &stats.Profit, &stats.UnitPrice, &stats.WholesalePrice)
		if err != nil {
			stats = nil
		}
	})
	return stats
}

// FindMany 查询多个财务统计信息
func (f *FinancialStatistics) FindMany() []*FinancialStatisticsAttrs {
	data := make([]*FinancialStatisticsAttrs, 0)
	f.Query(func(rows *sql.Rows) {
		tmp := new(FinancialStatisticsAttrs)
		_ = rows.Scan(&tmp.ID, &tmp.ProductID, &tmp.Profit, &tmp.UnitPrice, &tmp.WholesalePrice)
		data = append(data, tmp)
	})
	return data
}
