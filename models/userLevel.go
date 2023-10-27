package models

import (
	"database/sql"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

const (
	UserLevelStatusActivate = 10
	UserLevelStatusDisabled = -1
	UserLevelStatusDelete   = -2
)

type UserLevelAttrs struct {
	Id        int64   `json:"id"`         //主键
	AdminId   int64   `json:"admin_id"`   //管理员ID
	Name      string  `json:"name"`       //名称
	Icon      string  `json:"icon"`       //图标
	Level     int64   `json:"level"`      //等级
	Money     float64 `json:"money"`      //购买金额
	Days      int64   `json:"days"`       //购买天数 -1无限时间
	Status    int64   `json:"status"`     //状态 -2删除 -1禁用 10启用
	Data      string  `json:"data"`       //数据
	CreatedAt int64   `json:"created_at"` //创建时间
	UpdatedAt int64   `json:"updated_at"` //更新时间
}

type UserLevel struct {
	define.Db
}

func NewUserLevel(tx *sql.Tx) *UserLevel {
	return &UserLevel{
		database.DbPool.NewDb(tx).Table("user_level"),
	}
}

// AndWhere where条件
func (c *UserLevel) AndWhere(str string, arg ...any) *UserLevel {
	c.Db.AndWhere(str, arg...)
	return c
}

func (c *UserLevel) FindOne() *UserLevelAttrs {
	attrs := new(UserLevelAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.Name, &attrs.Icon, &attrs.Level, &attrs.Money, &attrs.Days, &attrs.Status, &attrs.Data, &attrs.CreatedAt, &attrs.UpdatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *UserLevel) FindMany() []*UserLevelAttrs {
	var data []*UserLevelAttrs
	c.Query(func(rows *sql.Rows) {
		tmp := new(UserLevelAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.Name, &tmp.Icon, &tmp.Level, &tmp.Money, &tmp.Days, &tmp.Status, &tmp.Data, &tmp.CreatedAt, &tmp.UpdatedAt)
		data = append(data, tmp)
	})
	return data
}
