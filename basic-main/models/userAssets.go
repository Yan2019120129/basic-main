package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

const (
	UserAssetsStatusActivate = 10
	UserAssetsStatusDelete   = -2
)

// UserAssetsAttrs 数据库模型属性
type UserAssetsAttrs struct {
	Id          int64   `json:"id"`           //主键
	AdminId     int64   `json:"admin_id"`     //管理员ID
	UserId      int64   `json:"user_id"`      //用户ID
	AssetsId    int64   `json:"assets_id"`    //资产ID
	Money       float64 `json:"money"`        //金额
	FreezeMoney float64 `json:"freeze_money"` //冻结金额
	Data        string  `json:"data"`         //数据
	Status      int64   `json:"status"`       //状态 -2删除｜-1禁用｜10启用
	CreatedAt   int64   `json:"created_at"`   //创建时间
	UpdatedAt   int64   `json:"updated_at"`   //更新时间
}

// UserAssets 数据库模型
type UserAssets struct {
	define.Db
}

// NewUserAssets 创建数据库模型
func NewUserAssets(tx *sql.Tx) *UserAssets {
	return &UserAssets{
		database.DbPool.NewDb(tx).Table("user_assets"),
	}
}

// AndWhere where条件
func (c *UserAssets) AndWhere(str string, arg ...any) *UserAssets {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *UserAssets) FindOne() *UserAssetsAttrs {
	attrs := new(UserAssetsAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.UserId, &attrs.AssetsId, &attrs.Money, &attrs.FreezeMoney, &attrs.Data, &attrs.Status, &attrs.CreatedAt, &attrs.UpdatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *UserAssets) FindMany() []*UserAssetsAttrs {
	data := make([]*UserAssetsAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(UserAssetsAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.UserId, &tmp.AssetsId, &tmp.Money, &tmp.FreezeMoney, &tmp.Data, &tmp.Status, &tmp.CreatedAt, &tmp.UpdatedAt)
		data = append(data, tmp)
	})
	return data
}

// 用户资产花费
func UserAssetsSpend(tx *sql.Tx, userInfo *UserAttrs, assetsInfo *AssetsAttrs, sourceId, billType int64, spendMoney float64) error {
	userAssetsInfo := NewUserAssets(tx).AndWhere("user_id=?", userInfo.Id).AndWhere("assets_id=?", assetsInfo.Id).AndWhere("status=?", UserAssetsStatusActivate).FindOne()
	if userAssetsInfo == nil || spendMoney <= 0 || userAssetsInfo.Money < spendMoney {
		return errors.New("insufficientBalance")
	}

	// 减去用户资产数量
	currentMoney := userAssetsInfo.Money - spendMoney
	_, err := NewUserAssets(tx).Value("money=?").Args(currentMoney).
		AndWhere("id=?", userAssetsInfo.Id).
		Update()
	if err != nil {
		return err
	}

	// 上级用户分销收益
	UserAssetsDistributorRevenue(tx, userInfo, assetsInfo, sourceId, billType, spendMoney)

	// 写入账单
	NewUserBill(tx).WriteUserBill(userInfo.AdminId, userInfo.Id, sourceId, billType, userAssetsInfo.Money, spendMoney)

	return nil
}

// 用户资产入金
func UserAssetsDeposit(tx *sql.Tx, userInfo *UserAttrs, assetsInfo *AssetsAttrs, sourceId, billType int64, depositMoney float64) error {
	if depositMoney <= 0 {
		return errors.New("IncorrectAmount")
	}

	userAssetsInfo := NewUserAssets(tx).AndWhere("user_id=?", userInfo.Id).AndWhere("assets_id=?", assetsInfo.Id).AndWhere("status=?", UserAssetsStatusActivate).FindOne()
	nowTime := time.Now()
	if userAssetsInfo == nil {
		//	新增用户资产
		_, err := NewUserAssets(tx).Field("admin_id", "user_id", "assets_id", "money", "created_at", "updated_at").
			Args(userInfo.AdminId, userInfo.Id, assetsInfo.Id, depositMoney, nowTime.Unix(), nowTime.Unix()).
			Insert()
		if err != nil {
			return err
		}
	} else {
		currentMoney := userAssetsInfo.Money + depositMoney
		_, err := NewUserAssets(tx).Value("money=?").Args(currentMoney).
			AndWhere("id=?", userAssetsInfo.Id).
			Update()
		if err != nil {
			return err
		}
	}

	// 上级用户分销收益
	UserAssetsDistributorRevenue(tx, userInfo, assetsInfo, sourceId, billType, depositMoney)

	// 写入账单
	NewUserBill(tx).WriteUserBill(userInfo.AdminId, userInfo.Id, sourceId, billType, userAssetsInfo.Money, depositMoney)

	return nil
}

// UserAssetsDistributorRevenue 用户资产分销商收益
func UserAssetsDistributorRevenue(tx *sql.Tx, userInfo *UserAttrs, assetsInfo *AssetsAttrs, sourceId, billType int64, involveMoney float64) error {
	if userInfo.ParentId <= 0 {
		return nil
	}

	settingAdminId := NewAdminUser(nil).GetSettingAdminId(userInfo.AdminId)

	//	判断是否开启账单类型收益
	revenueListStr := NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "pyramid_items")
	isRun, revenueBillType := GetRevenueBillType(billType, revenueListStr)
	if isRun {
		//	判断收益等级是否设置
		revenueLevelStr := NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "pyramid_level")
		var revenueLevel []map[string]float64
		_ = json.Unmarshal([]byte(revenueLevelStr), &revenueLevel)

		for i := 0; i < len(revenueLevel); i++ {
			if revenueLevel[i]["value"] <= 0 {
				continue
			}

			// 给上级收益
			revenueUserInfo := NewUser(nil).AndWhere("id=?", userInfo.ParentId).FindOne()
			revenueMoney := involveMoney * revenueLevel[i]["value"] / 100
			_ = UserAssetsDeposit(tx, revenueUserInfo, assetsInfo, sourceId, revenueBillType, revenueMoney)
		}
	}

	return nil
}
