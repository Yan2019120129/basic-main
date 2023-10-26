package models

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// CategoryAttrs 数据库模型属性
type CategoryAttrs struct {
	Id      int64  `json:"id"`       //类目ID
	AdminId string `json:"admin_id"` //类目图片
	Image   string `json:"image"`    //类目图片
	Name    string `json:"name"`     //种类名称
	Date    int64  `json:"date"`     //时间
	Status  int64  `json:"status"`   //状态 -2删除 -1禁用 10启用
	Operate int64  `json:"operate"`  //操作
}

const (
	CategoryActive   = 10 // 在使用
	CategoryDisabled = -1 // -1禁用
	CategoryDelete   = -2 // -2删除
)

// Category 数据库模型
type Category struct {
	define.Db
}

// NewCategory 创建数据库模型
func NewCategory(tx *sql.Tx) *Category {
	return &Category{
		database.DbPool.NewDb(tx).Table("category"),
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
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.Image, &attrs.Name, &attrs.Date, &attrs.Status, &attrs.Operate)
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
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.Image, &tmp.Name, &tmp.Date, &tmp.Status, &tmp.Operate)
		data = append(data, tmp)
	})
	return data
}
