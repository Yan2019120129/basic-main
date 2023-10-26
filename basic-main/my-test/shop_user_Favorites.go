package my_test

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// UserFavoritesStatus 表示用户收藏状态的常量
const (
	UserFavoritesStatusDeleted = -1 //已删除
	UserFavoritesStatusActive  = 1  //1活跃
	UserFavoritesStatusPending = 2  //2待定
)

// UserFavoritesAttrs 表示用户收藏信息的结构体
type UserFavoritesAttrs struct {
	ID          int64  `json:"id"`           // 收藏ID
	UserID      int64  `json:"user_id"`      // 用户ID
	CommodityID int64  `json:"commodity_id"` // 商品ID
	ProductName string `json:"product_name"` // 商品名称
	Date        int    `json:"date"`         // 收藏时间
	Status      int    `json:"status"`       // 收藏状态：-1已删除，1活跃，2待定
}

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

// AndWhere 添加 WHERE 子句
func (u *UserFavorites) AndWhere(str string, args ...interface{}) *UserFavorites {
	u.Db.AndWhere(str, args...)
	return u
}

// FindOne 查询单个收藏信息
func (u *UserFavorites) FindOne() *UserFavoritesAttrs {
	favorites := new(UserFavoritesAttrs)
	u.QueryRow(func(row *sql.Row) {
		err := row.Scan(&favorites.ID, &favorites.UserID, &favorites.CommodityID, &favorites.ProductName, &favorites.Date, &favorites.Status)
		if err != nil {
			favorites = nil
		}
	})
	return favorites
}

// FindMany 查询多个收藏信息
func (u *UserFavorites) FindMany() []*UserFavoritesAttrs {
	data := make([]*UserFavoritesAttrs, 0)
	u.Query(func(rows *sql.Rows) {
		tmp := new(UserFavoritesAttrs)
		_ = rows.Scan(&tmp.ID, &tmp.UserID, &tmp.CommodityID, &tmp.ProductName, &tmp.Date, &tmp.Status)
		data = append(data, tmp)
	})
	return data
}
