package models

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// UserFavoritesAttrs 数据库模型属性
type UserFavoritesAttrs struct {
	Id           int64  `json:"id"`           //收藏ID
	UserId       int64  `json:"user_id"`      //用户ID
	AdminId      int64  `json:"admin_id"`     //管理员ID
	CommodityId  int64  `json:"commodity_id"` //商品ID
	ProductName  string `json:"product_name"` //商品名称
	ProductImage string `json:"product_imag"` //商品图片
	Date         int64  `json:"date"`         //收藏时间
	Status       int64  `json:"status"`       //收藏状态：-1取消收藏，1收藏
}

const (
	UserFavoritesStatusActivate = 1  //	状态激活
	UserFavoritesStatusDelete   = -1 //	取消收藏
)

// UserFavorites 数据库模型
type UserFavorites struct {
	define.Db
}

// NewUserFavorites 创建数据库模型
func NewUserFavorites(tx *sql.Tx) *UserFavorites {
	return &UserFavorites{
		database.DbPool.NewDb(tx).Table("user_favorites"),
	}
}

// AndWhere where条件
func (c *UserFavorites) AndWhere(str string, arg ...any) *UserFavorites {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *UserFavorites) FindOne() *UserFavoritesAttrs {
	attrs := new(UserFavoritesAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.UserId, &attrs.AdminId, &attrs.CommodityId, &attrs.ProductName, &attrs.ProductImage, &attrs.Date, &attrs.Status)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *UserFavorites) FindMany() []*UserFavoritesAttrs {
	data := make([]*UserFavoritesAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(UserFavoritesAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.UserId, &tmp.AdminId, &tmp.CommodityId, &tmp.ProductName, &tmp.ProductImage, &tmp.Date, &tmp.Status)
		data = append(data, tmp)
	})
	return data
}
