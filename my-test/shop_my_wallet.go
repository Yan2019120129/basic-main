package my_test

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

type MyWalletAttrs struct {
	ID            int64   `json:"id"`             // 钱包唯一ID
	UserID        int64   `json:"user_id"`        // 用户ID
	Amount        float64 `json:"amount"`         // 金额
	PaymentMethod string  `json:"payment_method"` // 支付方式
	Date          int     `json:"date"`           // 交易时间
	Status        int     `json:"status"`         // 支付状态：-1待支付，1已支付
}

// MyWalletStatus 表示钱包支付状态
const (
	MyWalletStatusPending = -1 // 待支付
	MyWalletStatusPaid    = 1  // 已支付
)

// MyWallet 数据库模型
type MyWallet struct {
	define.Db
}

// NewMyWallet 创建数据库模型
func NewMyWallet(tx *sql.Tx) *MyWallet {
	return &MyWallet{
		database.DbPool.NewDb(tx).Table("my_wallet"),
	}
}

// AndWhere 添加 WHERE 子句
func (m *MyWallet) AndWhere(str string, args ...interface{}) *MyWallet {
	m.Db.AndWhere(str, args...)
	return m
}

// FindOne 查询单个钱包信息
func (m *MyWallet) FindOne() *MyWalletAttrs {
	wallet := new(MyWalletAttrs)
	m.QueryRow(func(row *sql.Row) {
		err := row.Scan(&wallet.ID, &wallet.UserID, &wallet.Amount, &wallet.PaymentMethod, &wallet.Date, &wallet.Status)
		if err != nil {
			wallet = nil
		}
	})
	return wallet
}

// FindMany 查询多个钱包信息
func (m *MyWallet) FindMany() []*MyWalletAttrs {
	data := make([]*MyWalletAttrs, 0)
	m.Query(func(rows *sql.Rows) {
		tmp := new(MyWalletAttrs)
		_ = rows.Scan(&tmp.ID, &tmp.UserID, &tmp.Amount, &tmp.PaymentMethod, &tmp.Date, &tmp.Status)
		data = append(data, tmp)
	})
	return data
}
