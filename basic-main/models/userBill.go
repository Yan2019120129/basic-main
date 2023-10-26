package models

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

const (
	UserBillTypeSystemDeposit         int64 = 1   //	系统充值
	UserBillTypeSystemDeduction       int64 = 2   //	系统扣除
	UserBillTypeDeposit               int64 = 3   //	用户充值
	UserBillTypeWithdraw              int64 = 4   //	用户提现
	UserBillTypeWithdrawRefuse        int64 = 5   //	提现拒绝
	UserBillTypeBuyLevel              int64 = 10  //	购买等级
	UserBillTypeBuyUpgradeLevel       int64 = 11  //	升级等级
	UserBillTypeRegisterRewards       int64 = 15  //	注册奖励
	UserBillTypeTaskRewards           int64 = 16  //	任务奖励
	UserBillTypeExtraRewards          int64 = 17  //	额外奖励
	UserBillTypeInviteRewards         int64 = 18  //	邀请奖励
	UserBillTypeBuyProduct            int64 = 20  //	购买产品
	UserBillTypeReturnProductAmount   int64 = 21  //	退回本金
	UserBillTypeProductProfit         int64 = 22  //	产品利润
	UserBillTypeFee                   int64 = 23  //	产品手续费
	UserBillTypeBuyProductEarnings    int64 = 30  //	分销购买产品收益
	UserBillTypeProductProfitEarnings int64 = 31  //	分销产品利润收益
	UserBillTypeAssetsSystemDeposit   int64 = 101 //	用户资产系统充值
	UserBillTypeAssetsSystemDeduction int64 = 102 //	用户资产系统扣除
)

// UserBillTypeNameMap 语言字典名称
var UserBillTypeNameMap = map[int64]string{
	UserBillTypeSystemDeposit: "systemDeposit", UserBillTypeSystemDeduction: "systemDeduction", UserBillTypeDeposit: "deposit",
	UserBillTypeWithdraw: "withdraw", UserBillTypeWithdrawRefuse: "withdrawRefuse", UserBillTypeBuyLevel: "buyLevel",
	UserBillTypeBuyUpgradeLevel: "buyUpgradeLevel", UserBillTypeRegisterRewards: "registerRewards", UserBillTypeTaskRewards: "taskRewards",
	UserBillTypeExtraRewards: "extraRewards", UserBillTypeInviteRewards: "inviteRewards", UserBillTypeBuyProduct: "buyProduct", UserBillTypeFee: "productFee", UserBillTypeReturnProductAmount: "returnProductAmount",
	UserBillTypeProductProfit: "productProfit", UserBillTypeBuyProductEarnings: "buyProductEarnings", UserBillTypeProductProfitEarnings: "productProfitEarnings",
	UserBillTypeAssetsSystemDeposit: "assetsSystemDeposit", UserBillTypeAssetsSystemDeduction: "assetsSystemDeduction",
}

type UserBillAttrs struct {
	Id        int64   `json:"id"`         //主键
	AdminId   int64   `json:"admin_id"`   //管理员ID
	UserId    int64   `json:"user_id"`    //用户ID
	SourceId  int64   `json:"source_id"`  //来源ID
	Name      string  `json:"name"`       //标题
	Type      int64   `json:"type"`       //类型
	Balance   float64 `json:"balance"`    //余额
	Money     float64 `json:"money"`      //金额
	Data      string  `json:"data"`       //数据
	CreatedAt int64   `json:"created_at"` //创建时间
}

type UserBill struct {
	define.Db
}

func NewUserBill(tx *sql.Tx) *UserBill {
	return &UserBill{
		database.DbPool.NewDb(tx).Table("user_bill"),
	}
}

// AndWhere where条件
func (c *UserBill) AndWhere(str string, arg ...any) *UserBill {
	c.Db.AndWhere(str, arg...)
	return c
}

func (c *UserBill) FindOne() *UserBillAttrs {
	attrs := new(UserBillAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.UserId, &attrs.SourceId, &attrs.Name, &attrs.Type, &attrs.Balance, &attrs.Money, &attrs.Data, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// WriteUserBill 写入账单
func (c *UserBill) WriteUserBill(adminId, userId, sourceId, billType int64, beforeMoney, money float64) {
	nowTime := time.Now()
	_, err := c.Field("admin_id", "user_id", "source_id", "name", "type", "balance", "money", "created_at").
		Args(adminId, userId, sourceId, UserBillTypeNameMap[billType], billType, beforeMoney, money, nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}
}

// GetBillTypeMoney 获取账单类型金额
func GetBillTypeMoney(billType int64, beforeMoney, money float64) float64 {
	if billType == UserBillTypeSystemDeduction || billType == UserBillTypeAssetsSystemDeduction || billType == UserBillTypeWithdraw || billType == UserBillTypeBuyLevel ||
		billType == UserBillTypeBuyUpgradeLevel || billType == UserBillTypeBuyProduct {
		return beforeMoney - money
	}
	return beforeMoney + money
}

// GetRevenueBillType 获取收益的类型
func GetRevenueBillType(billType int64, itemsStr string) (bool, int64) {
	var revenueList map[string]bool
	_ = json.Unmarshal([]byte(itemsStr), &revenueList)

	for k, isTrue := range revenueList {
		itemList := strings.Split(k, "_")
		if len(itemList) == 2 && isTrue {
			currentBillType, _ := strconv.ParseInt(itemList[0], 10, 64)
			revenueBillType, _ := strconv.ParseInt(itemList[0], 10, 64)
			if billType == currentBillType {
				return true, revenueBillType
			}
		}
	}

	return false, 0
}
