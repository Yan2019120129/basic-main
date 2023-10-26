package models

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// FinancialStatisticsAttrs 数据库模型属性
type FinancialStatisticsAttrs struct {
	Id             int64   `json:"id"`              //财务统计ID
	AdminId        int64   `json:"admin_id"`        //管理员ID
	ProductId      int64   `json:"product_id"`      //商品ID
	Profit         float64 `json:"profit"`          //利润
	UnitPrice      float64 `json:"unit_price"`      //单价
	WholesalePrice float64 `json:"wholesale_price"` //批发价
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

// AndWhere where条件
func (c *FinancialStatistics) AndWhere(str string, arg ...any) *FinancialStatistics {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *FinancialStatistics) FindOne() *FinancialStatisticsAttrs {
	attrs := new(FinancialStatisticsAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.ProductId, &attrs.Profit, &attrs.UnitPrice, &attrs.WholesalePrice)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *FinancialStatistics) FindMany() []*FinancialStatisticsAttrs {
	data := make([]*FinancialStatisticsAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(FinancialStatisticsAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.ProductId, &tmp.Profit, &tmp.UnitPrice, &tmp.WholesalePrice)
		data = append(data, tmp)
	})
	return data
}
