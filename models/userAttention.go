package models

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// UserAttentionAttrs 数据库模型属性
type UserAttentionAttrs struct {
	Id       int64  `json:"id"`        //关注ID
	UserId   int64  `json:"user_id"`   //用户ID
	AdminId  int64  `json:"admin_id"`  //管理员ID
	ShopId   int64  `json:"shop_id"`   //店铺ID
	ShopLogo string `json:"shop_logo"` //店铺log
	ShopName string `json:"shop_name"` //店铺名
	Date     int64  `json:"date"`      //时间
	Status   int64  `json:"status"`    //状态：1关注，2取消关注
}

const (
	UserAttentionActivate = 1  //	关注
	UserAttentionDelete   = -1 //	取消关注
)

// UserAttention 数据库模型
type UserAttention struct {
	define.Db
}

// NewUserAttention 创建数据库模型
func NewUserAttention(tx *sql.Tx) *UserAttention {
	return &UserAttention{
		database.DbPool.NewDb(tx).Table("user_attention"),
	}
}

// AndWhere where条件
func (c *UserAttention) AndWhere(str string, arg ...any) *UserAttention {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *UserAttention) FindOne() *UserAttentionAttrs {
	attrs := new(UserAttentionAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.UserId, &attrs.AdminId, &attrs.ShopId, &attrs.ShopLogo, &attrs.ShopName, &attrs.Date, &attrs.Status)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *UserAttention) FindMany() []*UserAttentionAttrs {
	data := make([]*UserAttentionAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(UserAttentionAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.UserId, &tmp.AdminId, &tmp.ShopId, &tmp.ShopLogo, &tmp.ShopName, &tmp.Date, &tmp.Status)
		data = append(data, tmp)
	})
	return data
}
