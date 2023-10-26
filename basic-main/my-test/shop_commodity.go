package my_test

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// CommodityAttrs 表示会话信息的结构体
type CommodityAttrs struct {
	ID             int    `json:"id"`              // 主键
	ConversationID string `json:"conversation_id"` // 会话ID
	UserID         int    `json:"user_id"`         // 用户ID
	ReceiverID     int    `json:"receiver_id"`     // 接收者用户ID
	Type           int    `json:"type"`            // 类型，1表示在线客服
	Status         int    `json:"status"`          // 状态，-2表示删除，-1表示屏蔽，10表示正常
	Data           string `json:"data"`            // 数据
	UpdatedAt      int    `json:"updated_at"`      // 更新时间
	CreatedAt      int    `json:"created_at"`      // 创建时间
}

const (
	CommodityStatusDelete = -2
	CommodityStatusShield = -1
	CommodityStatusNormal = 10
)

// Commodity 数据库模型
type Commodity struct {
	define.Db
}

// NewCommodity 创建数据库模型
func NewCommodity(tx *sql.Tx) *Commodity {
	return &Commodity{
		database.DbPool.NewDb(tx).Table("commodity"),
	}
}

// AndWhere 添加 WHERE 子句
func (c *Commodity) AndWhere(str string, args ...interface{}) *Commodity {
	c.Db.AndWhere(str, args...)
	return c
}

// FindOne 查询单个商品信息
func (c *Commodity) FindOne() *CommodityAttrs {
	commodity := new(CommodityAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&commodity.ID, &commodity.ConversationID, &commodity.UserID, &commodity.ReceiverID, &commodity.Type, &commodity.Status, &commodity.Data, &commodity.UpdatedAt, &commodity.CreatedAt)
		if err != nil {
			commodity = nil
		}
	})
	return commodity
}

// FindMany 查询多个商品信息
func (c *Commodity) FindMany() []*CommodityAttrs {
	data := make([]*CommodityAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(CommodityAttrs)
		_ = rows.Scan(&tmp.ID, &tmp.ConversationID, &tmp.UserID, &tmp.ReceiverID, &tmp.Type, &tmp.Status, &tmp.Data, &tmp.UpdatedAt, &tmp.CreatedAt)
		data = append(data, tmp)
	})
	return data
}
