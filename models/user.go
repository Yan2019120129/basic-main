package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

const (
	UserStatusDelete   = -2 //	用户删除状态
	UserStatusDisabled = -1 //	用户禁用状态
	UserStatusFreeze   = 1  //	用户冻结状态
	UserStatusActivate = 10 //	激活用户状态

	UserTypeOnline    = -10 //	客服坐席
	UserTypeTemporary = -2  //	临时用户
	UserTypeVirtual   = -1  //	虚拟用户
	UserTypeReality   = 10  //	真实用户
	UserTypeVip       = 20  //	VIP用户
)

// UserAttrs 数据库模型属性
type UserAttrs struct {
	Id          int64   `json:"id"`           //主键
	AdminId     int64   `json:"admin_id"`     //管理员ID
	ParentId    int64   `json:"parent_id"`    //上级ID
	CountryId   int64   `json:"country_id"`   //国家ID
	UserName    string  `json:"username"`     //用户名
	Nickname    string  `json:"nickname"`     //昵称
	Email       string  `json:"email"`        //邮箱
	Telephone   string  `json:"telephone"`    //手机号码
	Avatar      string  `json:"avatar"`       //头像
	Sex         int64   `json:"sex"`          //类型 -1未知 1男 2女
	Birthday    int64   `json:"birthday"`     //生日
	Password    string  `json:"password"`     //密码
	SecurityKey string  `json:"security_key"` //安全密钥
	Money       float64 `json:"money"`        //金额
	FreezeMoney float64 `json:"freeze_money"` //冻结金额
	Type        int64   `json:"type"`         //类型 -2临时用户 -1虚拟 10普通 20会员用户
	Status      int64   `json:"status"`       //状态 -2删除｜-1禁用｜10启用
	Data        string  `json:"data"`         //数据
	Ip4         string  `json:"ip4"`          //IP4地址
	CreatedAt   int64   `json:"created_at"`   //创建时间
	UpdatedAt   int64   `json:"updated_at"`   //更新时间
}

type UserVerifyInfo struct {
	Status int64  `json:"status"` //	验证状态
	Data   string `json:"data"`   //	信息
}

type UserLevelInfo struct {
	Id        int64   `json:"id"`        //	Id
	Level     int64   `json:"level"`     //	等级
	Name      string  `json:"name"`      //	名称
	Icon      string  `json:"icon"`      //	图标
	Days      int64   `json:"days"`      //	天数
	Money     float64 `json:"money"`     //	金额
	CreatedAt int64   `json:"createdAt"` //	创建时间戳
	UpdatedAt int64   `json:"updatedAt"` //	过期时间戳
}

type UserInfo struct {
	Id          int64           `json:"id"`           //	用户ID
	CountryId   int64           `json:"countryId"`    //	国家ID
	UserName    string          `json:"username"`     //	用户名
	Nickname    string          `json:"nickname"`     //	昵称
	Email       string          `json:"email"`        //	邮箱
	Telephone   string          `json:"telephone"`    //	手机号
	Avatar      string          `json:"avatar"`       //	头像
	Sex         int64           `json:"sex"`          //	类型 -1未知 1男 2女
	Birthday    int64           `json:"birthday"`     //	生日
	Money       float64         `json:"money"`        //	金额
	FreezeMoney float64         `json:"freeze_money"` //	冻结金额
	Data        string          `json:"data"`         //	数据
	VerifyInfo  *UserVerifyInfo `json:"verifyInfo"`   //	是否验证 -1验证失败 0未验证 10 正在验证 20已验证
	LevelInfo   *UserLevelInfo  `json:"levelInfo"`    //	等级信息
	InviteCode  string          `json:"inviteCode"`   //	邀请码
	CreatedAt   int64           `json:"createdAt"`    //	创建时间
	UpdatedAt   int64           `json:"updatedAt"`    //	更新时间
}

// UserTree 用户树
type UserTree struct {
	Id          int64       `json:"id"`           //	用户ID
	Header      string      `json:"header"`       //	层级
	Avatar      string      `json:"avatar"`       //	头像
	UserName    string      `json:"username"`     //	用户名
	SumPeople   int64       `json:"sum_people"`   //	总人数
	SumAmount   float64     `json:"sum_amount"`   //	总充值
	SumEarnings float64     `json:"sum_earnings"` //	总收益
	Children    []*UserTree `json:"children"`     //	子集
}

// User 数据库模型
type User struct {
	define.Db
}

// NewUser 创建数据库模型
func NewUser(tx *sql.Tx) *User {
	return &User{
		database.DbPool.NewDb(tx).Table("user"),
	}
}

// AndWhere where条件
func (c *User) AndWhere(str string, arg ...any) *User {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *User) FindOne() *UserAttrs {
	attrs := new(UserAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.ParentId, &attrs.CountryId, &attrs.UserName, &attrs.Nickname, &attrs.Email, &attrs.Telephone, &attrs.Avatar, &attrs.Sex, &attrs.Birthday, &attrs.Password, &attrs.SecurityKey, &attrs.Money, &attrs.FreezeMoney, &attrs.Type, &attrs.Status, &attrs.Data, &attrs.Ip4, &attrs.CreatedAt, &attrs.UpdatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *User) FindMany() []*UserAttrs {
	data := make([]*UserAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(UserAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.ParentId, &tmp.CountryId, &tmp.UserName, &tmp.Nickname, &tmp.Email, &tmp.Telephone, &tmp.Avatar, &tmp.Sex, &tmp.Birthday, &tmp.Password, &tmp.SecurityKey, &tmp.Money, &tmp.FreezeMoney, &tmp.Type, &tmp.Status, &tmp.Data, &tmp.Ip4, &tmp.CreatedAt, &tmp.UpdatedAt)
		data = append(data, tmp)
	})
	return data
}

// FindUserLikeNameIds 获取用户名称IDS
func (c *User) FindUserLikeNameIds(username string) []string {
	data := NewUser(nil).
		Field("id").
		AndWhere("username like ?", "%"+username+"%").ColumnString()
	if len(data) == 0 {
		return []string{"-1"}
	}
	return data
}

// GetUserTree 获取用户树
func (c *User) GetUserTree(adminId int64, parentId int64) []*UserTree {
	adminIds := NewAdminUser(nil).GetAdminChildrenParentIds(adminId)

	data := make([]*UserTree, 0)
	NewUser(nil).Field("id", "avatar", "username").
		AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").
		AndWhere("status>?", UserStatusDelete).AndWhere("parent_id=?", parentId).Query(func(rows *sql.Rows) {
		temp := new(UserTree)
		_ = rows.Scan(&temp.Id, &temp.Avatar, &temp.UserName)
		if parentId == 0 {
			temp.Header = "root"
		} else {
			temp.Header = "generic"
		}

		// 总人数
		NewUser(nil).Field("count(*)").
			AndWhere("parent_id=?", temp.Id).
			AndWhere("status>?", UserStatusDelete).QueryRow(func(row *sql.Row) {
			_ = row.Scan(&temp.SumPeople)
		})

		// 总充值
		NewUserWalletOrder(nil).Field("sum(money)").AndWhere("user_id=?", temp.Id).AndWhere("type in (" + strconv.FormatInt(WalletOrderTypeDeposit, 10) + "," + strconv.FormatInt(WalletOrderTypeSystemDeposit, 10) + ")").QueryRow(func(row *sql.Row) {
			_ = row.Scan(&temp.SumAmount)
		})

		// 总收益
		NewUserBill(nil).Field("sum(money)").AndWhere("user_id=?", temp.Id).AndWhere("type=?", UserBillTypeProductProfit).QueryRow(func(row *sql.Row) {
			_ = row.Scan(&temp.SumEarnings)
		})

		// 子集用户
		temp.Children = NewUser(nil).GetUserTree(adminId, temp.Id)
		data = append(data, temp)
	})

	return data
}

// UserSpend 用户花费
func UserSpend(tx *sql.Tx, userInfo *UserAttrs, sourceId, billType int64, spendMoney float64) error {
	if spendMoney <= 0 || userInfo.Money < spendMoney {
		return errors.New("insufficientBalance")
	}

	// 更新用户信息
	currentMoney := userInfo.Money - spendMoney
	_, err := NewUser(tx).Value("money=?").Args(currentMoney).AndWhere("id=?", userInfo.Id).Update()
	if err != nil {
		return err
	}

	// 上级用户分销收益
	UserDistributorRevenue(tx, userInfo, sourceId, billType, spendMoney)

	// 写入账单
	NewUserBill(tx).WriteUserBill(userInfo.AdminId, userInfo.Id, sourceId, billType, userInfo.Money, spendMoney)
	return nil
}

// UserDeposit 用户入金
func UserDeposit(tx *sql.Tx, userInfo *UserAttrs, sourceId, billType int64, depositMoney float64) error {
	if depositMoney <= 0 {
		return errors.New("IncorrectAmount")
	}

	// 更新用户信息
	currentMoney := userInfo.Money + depositMoney
	_, err := NewUser(tx).Value("money=?").Args(currentMoney).AndWhere("id=?", userInfo.Id).Update()
	if err != nil {
		return err
	}

	// 上级用户分销收益
	UserDistributorRevenue(tx, userInfo, sourceId, billType, depositMoney)

	// 写入账单
	NewUserBill(tx).WriteUserBill(userInfo.AdminId, userInfo.Id, sourceId, billType, userInfo.Money, depositMoney)
	return nil
}

// UserDistributorRevenue 用户分销商收益
func UserDistributorRevenue(tx *sql.Tx, userInfo *UserAttrs, sourceId, billType int64, involveMoney float64) error {
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
			_ = UserDeposit(tx, revenueUserInfo, sourceId, revenueBillType, revenueMoney)
		}
	}

	return nil
}
