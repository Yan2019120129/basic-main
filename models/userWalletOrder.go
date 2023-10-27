package models

import (
	"database/sql"
	"time"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

const (
	// WalletOrderTypeDeposit 充值
	WalletOrderTypeDeposit = 1
	// WalletOrderTypeSystemDeposit 系统加款
	WalletOrderTypeSystemDeposit = 2
	// WalletOrderTypeWithdraw 提现
	WalletOrderTypeWithdraw = 10
	// WalletOrderTypeSystemWithdraw 系统减款
	WalletOrderTypeSystemWithdraw = 11
	//	WalletOrderTypeAssetsWithdraw 资产提现
	WalletOrderTypeAssetsWithdraw = 12

	WalletOrderStatusPending  = 10
	WalletOrderStatusRefuse   = -1
	WalletOrderStatusComplete = 20
	WalletOrderStatusDelete   = -2
)

type UserWalletOrderAttrs struct {
	Id        int64   `json:"id"`         //主键
	OrderSn   string  `json:"order_sn"`   //订单号
	AdminId   int64   `json:"admin_id"`   //管理员ID
	UserId    int64   `json:"user_id"`    //用户ID
	UserType  int64   `json:"user_type"`  //用户类型
	Type      int64   `json:"type"`       //类型 1充值 2系统加款 10提现 11系统减款
	PaymentId int64   `json:"payment_id"` //充值｜提现ID
	Money     float64 `json:"money"`      //金额
	Balance   float64 `json:"balance"`    //余额
	Status    int64   `json:"status"`     //状态 -1拒绝 10处理 20完成
	Proof     string  `json:"proof"`      //凭证
	Data      string  `json:"data"`       //数据
	Fee       float64 `json:"fee"`        //手续费
	UpdatedAt int64   `json:"updated_at"` //更新时间
	CreatedAt int64   `json:"created_at"` //创建时间
}

type UserWalletOrder struct {
	define.Db
}

func NewUserWalletOrder(tx *sql.Tx) *UserWalletOrder {
	return &UserWalletOrder{
		database.DbPool.NewDb(tx).Table("user_wallet_order"),
	}
}

// AndWhere where条件
func (c *UserWalletOrder) AndWhere(str string, arg ...any) *UserWalletOrder {
	c.Db.AndWhere(str, arg...)
	return c
}

func (c *UserWalletOrder) FindOne() *UserWalletOrderAttrs {
	attrs := new(UserWalletOrderAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.OrderSn, &attrs.AdminId, &attrs.UserId, &attrs.UserType, &attrs.Type, &attrs.PaymentId, &attrs.Money, &attrs.Balance, &attrs.Status, &attrs.Proof, &attrs.Data, &attrs.Fee, &attrs.UpdatedAt, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// WithdrawNums 用户时间范围次数
func (c *UserWalletOrder) WithdrawNums(userId int64, beforeTime time.Time) int64 {
	return c.Field("count(*)").
		AndWhere("user_id=?", "", userId).
		AndWhere("created_at>?", beforeTime.Unix()).Count()
}
