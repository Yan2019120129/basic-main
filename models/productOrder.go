package models

import (
	"database/sql"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// ProductOrderAttrs 数据库模型属性
type ProductOrderAttrs struct {
	Id        int64   `json:"id"`         //主键
	AdminId   int64   `json:"admin_id"`   //管理员ID
	UserId    int64   `json:"user_id"`    //用户ID
	ProductId int64   `json:"product_id"` //产品ID
	OrderSn   string  `json:"order_sn"`   //订单编号
	Money     float64 `json:"money"`      //金额
	Nums      int64   `json:"nums"`       //数量
	Type      int64   `json:"type"`       //类型
	Status    int64   `json:"status"`     //状态 -2删除 -1完结 10启用
	Data      string  `json:"data"`       //数据
	ExpiredAt int64   `json:"expired_at"` //过期时间
	UpdatedAt int64   `json:"updated_at"` //更新时间
	CreatedAt int64   `json:"created_at"` //创建时间
}

const (
	ProductOrderStatusComplete = -1
	ProductOrderStatusPending  = 10
)

// ProductOrder 数据库模型
type ProductOrder struct {
	define.Db
}

// NewProductOrder 创建数据库模型
func NewProductOrder(tx *sql.Tx) *ProductOrder {
	return &ProductOrder{
		database.DbPool.NewDb(tx).Table("product_order"),
	}
}

// AndWhere where条件
func (c *ProductOrder) AndWhere(str string, arg ...any) *ProductOrder {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *ProductOrder) FindOne() *ProductOrderAttrs {
	attrs := new(ProductOrderAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.UserId, &attrs.ProductId, &attrs.OrderSn, &attrs.Money, &attrs.Nums, &attrs.Type, &attrs.Status, &attrs.Data, &attrs.ExpiredAt, &attrs.UpdatedAt, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *ProductOrder) FindMany() []*ProductOrderAttrs {
	data := make([]*ProductOrderAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(ProductOrderAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.UserId, &tmp.ProductId, &tmp.OrderSn, &tmp.Money, &tmp.Nums, &tmp.Type, &tmp.Status, &tmp.Data, &tmp.ExpiredAt, &tmp.UpdatedAt, &tmp.CreatedAt)
		data = append(data, tmp)
	})
	return data
}
