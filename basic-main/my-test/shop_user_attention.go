package my_test

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// AttentionAttrs 表示关注信息的结构体
type AttentionAttrs struct {
	ID       int64  `json:"id"`        // 关注ID
	UserID   int64  `json:"user_id"`   // 用户ID
	ShopID   int64  `json:"shop_id"`   // 店铺ID
	ShopLogo string `json:"shop_logo"` // 店铺log
	ShopName string `json:"shop_name"` // 店铺名
	Date     int    `json:"date"`      // 时间
	Status   int    `json:"status"`    // 状态：-1，1，2
}

// Attention 数据库模型
type Attention struct {
	define.Db
}

// NewAttention 创建数据库模型
func NewAttention(tx *sql.Tx) *Attention {
	return &Attention{
		database.DbPool.NewDb(tx).Table("attention"),
	}
}

// AndWhere 添加 WHERE 子句
func (a *Attention) AndWhere(str string, args ...interface{}) *Attention {
	a.Db.AndWhere(str, args...)
	return a
}

// FindOne 查询单个关注信息
func (a *Attention) FindOne() *AttentionAttrs {
	attention := new(AttentionAttrs)
	a.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attention.ID, &attention.UserID, &attention.ShopID, &attention.ShopLogo, &attention.ShopName, &attention.Date, &attention.Status)
		if err != nil {
			attention = nil
		}
	})
	return attention
}

// FindMany 查询多个关注信息
func (a *Attention) FindMany() []*AttentionAttrs {
	data := make([]*AttentionAttrs, 0)
	a.Query(func(rows *sql.Rows) {
		tmp := new(AttentionAttrs)
		_ = rows.Scan(&tmp.ID, &tmp.UserID, &tmp.ShopID, &tmp.ShopLogo, &tmp.ShopName, &tmp.Date, &tmp.Status)
		data = append(data, tmp)
	})
	return data
}
