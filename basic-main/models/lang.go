package models

import (
	"database/sql"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

const (
	LangStatusActivate = 10
	LangStatusDisabled = -1
)

// LangAttrs 数据库模型属性
type LangAttrs struct {
	Id        int64  `json:"id"`         //主键
	AdminId   int64  `json:"admin_id"`   //管理员ID
	Name      string `json:"name"`       //名称
	Alias     string `json:"alias"`      //别名
	Icon      string `json:"icon"`       //图标
	Sort      int64  `json:"sort"`       //排序
	Status    int64  `json:"status"`     //状态 1禁用｜10启用
	Data      string `json:"data"`       //数据
	CreatedAt int64  `json:"created_at"` //创建时间
}

// Lang 数据库模型
type Lang struct {
	define.Db
}

// NewLang 创建数据库模型
func NewLang(tx *sql.Tx) *Lang {
	return &Lang{
		database.DbPool.NewDb(tx).Table("lang"),
	}
}

// AndWhere where条件
func (c *Lang) AndWhere(str string, arg ...any) *Lang {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *Lang) FindOne() *LangAttrs {
	attrs := new(LangAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.Name, &attrs.Alias, &attrs.Icon, &attrs.Sort, &attrs.Status, &attrs.Data, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *Lang) FindMany() []*LangAttrs {
	var data []*LangAttrs
	c.Query(func(rows *sql.Rows) {
		tmp := new(LangAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.Name, &tmp.Alias, &tmp.Icon, &tmp.Sort, &tmp.Status, &tmp.Data, &tmp.CreatedAt)
		data = append(data, tmp)
	})
	return data
}

// FindAlias 查询语言别名
func (c *Lang) FindAlias(adminId int64, alias string) *LangAttrs {
	c.AndWhere("admin_id=?", adminId).AndWhere("alias=?", alias)
	return c.FindOne()
}
