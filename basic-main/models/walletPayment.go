package models

import (
	"database/sql"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

const (
	WalletPaymentModeDeposit        = 1
	WalletPaymentModeWithdraw       = 10
	WalletPaymentTypeBank           = 1
	WalletPaymentTypeCryptocurrency = 10
	WalletPaymentTypeExternal       = 20
	WalletPaymentStatusActivate     = 10
	WalletPaymentStatusDelete       = -2
)

type WalletPaymentAttrs struct {
	Id          int64  `json:"id"`           //主键
	AdminId     int64  `json:"admin_id"`     //管理员ID
	Icon        string `json:"icon"`         //图标
	Mode        int64  `json:"mode"`         //方式 1充值 10提现
	Type        int64  `json:"type"`         //类型 1银行转账 10数字货币 20三方支付
	Name        string `json:"name"`         //名称
	AccountName string `json:"account_name"` //张三｜ERC20
	AccountCode string `json:"account_code"` //卡号｜地址
	Sort        int64  `json:"sort"`         //排序
	Status      int64  `json:"status"`       //状态 -2删除 -1禁用 10启用
	Description string `json:"description"`  //描述
	Data        string `json:"data"`         //数据
	Expand      string `json:"expand"`       //扩展数据
	CreatedAt   int64  `json:"created_at"`   //创建时间
	UpdatedAt   int64  `json:"updated_at"`   //更新时间
}

type WalletPayment struct {
	define.Db
}

func NewWalletPayment(tx *sql.Tx) *WalletPayment {
	return &WalletPayment{
		database.DbPool.NewDb(tx).Table("wallet_payment"),
	}
}

// AndWhere where条件
func (c *WalletPayment) AndWhere(str string, arg ...any) *WalletPayment {
	c.Db.AndWhere(str, arg...)
	return c
}

func (c *WalletPayment) FindOne() *WalletPaymentAttrs {
	attrs := new(WalletPaymentAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.Icon, &attrs.Mode, &attrs.Type, &attrs.Name, &attrs.AccountName, &attrs.AccountCode, &attrs.Sort, &attrs.Status, &attrs.Description, &attrs.Data, &attrs.Expand, &attrs.CreatedAt, &attrs.UpdatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *WalletPayment) FindMany() []*WalletPaymentAttrs {
	var data []*WalletPaymentAttrs
	c.Query(func(rows *sql.Rows) {
		tmp := new(WalletPaymentAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.Icon, &tmp.Mode, &tmp.Type, &tmp.Name, &tmp.AccountName, &tmp.AccountCode, &tmp.Sort, &tmp.Status, &tmp.Description, &tmp.Data, &tmp.Expand, &tmp.CreatedAt, &tmp.UpdatedAt)
		data = append(data, tmp)
	})
	return data
}
