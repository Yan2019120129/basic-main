package models

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// ShopUserAttrs 数据库模型属性
type ShopUserAttrs struct {
	Id	int64	`json:"id"`	//用户唯一Id
	HeadPortrait	string	`json:"head_portrait"`	//用户头像
	Name	string	`json:"name"`	//用户姓名
	Sex	int64	`json:"sex"`	//用户性别：-1未知，1男，2女
	Birthday	int64	`json:"birthday"`	//用户生日（存储时间戳）
	Email	string	`json:"email"`	//用户邮箱
	Phone	string	`json:"phone"`	//用户电话
	UpdateTime	int64	`json:"update_time"`	//用户信息最新修改时间（时间戳格式）
	CreateTime	int64	`json:"create_time"`	//用户注册时间(时间戳格式)
	Status	int64	`json:"status"`	//用户状态：-2删除，-1禁用，10启用
	Balance	float64	`json:"balance"`	//用户余额
}

// ShopUser 数据库模型
type ShopUser struct {
	define.Db
}

// NewShopUser 创建数据库模型
func NewShopUser(tx *sql.Tx) *ShopUser {
	return &ShopUser{
		database.DbPool.NewDb(tx).Table("shop_user"),
	}
}

// AndWhere where条件
func (c *ShopUser) AndWhere(str string, arg ...any) *ShopUser {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *ShopUser) FindOne() *ShopUserAttrs {
	attrs := new(ShopUserAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.HeadPortrait, &attrs.Name, &attrs.Sex, &attrs.Birthday, &attrs.Email, &attrs.Phone, &attrs.UpdateTime, &attrs.CreateTime, &attrs.Status, &attrs.Balance)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *ShopUser) FindMany() []*ShopUserAttrs {
	data := make([]*ShopUserAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(ShopUserAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.HeadPortrait, &tmp.Name, &tmp.Sex, &tmp.Birthday, &tmp.Email, &tmp.Phone, &tmp.UpdateTime, &tmp.CreateTime, &tmp.Status, &tmp.Balance)
		data = append(data, tmp)
	})
	return data
}