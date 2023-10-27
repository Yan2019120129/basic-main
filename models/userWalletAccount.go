package models

import (
	"database/sql"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

const (
	UserWalletAccountStatusActivate = 10
	UserWalletAccountStatusDelete   = -2
	UserWalletAccountStatusDisabled = -1
)

type UserWalletAccountAttrs struct {
	Id         int64  `json:"id"`          //主键
	AdminId    int64  `json:"admin_id"`    //管理员ID
	UserId     int64  `json:"user_id"`     //用户ID
	PaymentId  int64  `json:"payment_id"`  //提现方式ID
	Name       string `json:"name"`        //名称
	RealName   string `json:"real_name"`   //真实名字
	CardNumber string `json:"card_number"` //卡号
	Address    string `json:"address"`     //地址
	Status     int64  `json:"status"`      //状态 -2删除 -1禁用 10启用
	Sort       int64  `json:"sort"`        //排序
	Data       string `json:"data"`        //数据
	UpdatedAt  int64  `json:"updated_at"`  //更新时间
	CreatedAt  int64  `json:"created_at"`  //创建时间
}

type UserWalletAccount struct {
	define.Db
}

func NewUserWalletAccount(tx *sql.Tx) *UserWalletAccount {
	return &UserWalletAccount{
		database.DbPool.NewDb(tx).Table("user_wallet_account"),
	}
}

// AndWhere where条件
func (c *UserWalletAccount) AndWhere(str string, arg ...any) *UserWalletAccount {
	c.Db.AndWhere(str, arg...)
	return c
}

func (c *UserWalletAccount) FindOne() *UserWalletAccountAttrs {
	attrs := new(UserWalletAccountAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.UserId, &attrs.PaymentId, &attrs.Name, &attrs.RealName, &attrs.CardNumber, &attrs.Address, &attrs.Status, &attrs.Sort, &attrs.Data, &attrs.UpdatedAt, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// AccountNums 用户绑定类型数量
func (c *UserWalletAccount) AccountNums(userId int64, paymentId int64) int64 {
	return c.Field("count(*)").
		AndWhere("user_id=?", userId).
		AndWhere("status=?", UserWalletAccountStatusActivate).
		AndWhere("payment_id=?", paymentId).Count()
}
