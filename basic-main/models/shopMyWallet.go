package models

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// ShopMyWalletAttrs 数据库模型属性
type ShopMyWalletAttrs struct {
	Id	int64	`json:"id"`	//钱包唯一ID
	UserId	int64	`json:"user_id"`	//用户ID
	Amount	float64	`json:"amount"`	//金额
	PaymentMethod	string	`json:"payment_method"`	//支付方式
	Date	int64	`json:"date"`	//交易时间
	Status	int64	`json:"status"`	//支付状态：-1待支付，1已支付
}

// ShopMyWallet 数据库模型
type ShopMyWallet struct {
	define.Db
}

// NewShopMyWallet 创建数据库模型
func NewShopMyWallet(tx *sql.Tx) *ShopMyWallet {
	return &ShopMyWallet{
		database.DbPool.NewDb(tx).Table("shop_my_wallet"),
	}
}

// AndWhere where条件
func (c *ShopMyWallet) AndWhere(str string, arg ...any) *ShopMyWallet {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *ShopMyWallet) FindOne() *ShopMyWalletAttrs {
	attrs := new(ShopMyWalletAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.UserId, &attrs.Amount, &attrs.PaymentMethod, &attrs.Date, &attrs.Status)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *ShopMyWallet) FindMany() []*ShopMyWalletAttrs {
	data := make([]*ShopMyWalletAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(ShopMyWalletAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.UserId, &tmp.Amount, &tmp.PaymentMethod, &tmp.Date, &tmp.Status)
		data = append(data, tmp)
	})
	return data
}