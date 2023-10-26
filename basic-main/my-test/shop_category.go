package my_test

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

type CategoryAttrs struct {
	Id      int64  `json:"id"`      // 类目ID
	Image   string `json:"image"`   // 类目图片
	Name    string `json:"name"`    // 种类名称
	Date    int    `json:"date"`    // 时间
	Status  int    `json:"status"`  // 状态
	Operate int    `json:"operate"` // 操作
}

const (
	CategoryStatus    = -1
	CategoryStatusTwo = 10
)

// Category 数据库模型
type Category struct {
	define.Db
}

// NewCategory 创建数据库模型
func NewCategory(tx *sql.Tx) *Category {
	return &Category{
		database.DbPool.NewDb(tx).Table("product_order"),
	}
}

// AndWhere where条件
func (c *Category) AndWhere(str string, arg ...any) *Category {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *Category) FindOne() *CategoryAttrs {
	attrs := new(CategoryAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.Id, &attrs.Image, &attrs.Name, &attrs.Date, &attrs.Status, &attrs.Operate)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *Category) FindMany() []*CategoryAttrs {
	data := make([]*CategoryAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(CategoryAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.Id, &tmp.Image, &tmp.Name, &tmp.Date, &tmp.Status, &tmp.Operate)
		data = append(data, tmp)
	})
	return data
}
