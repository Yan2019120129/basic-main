package models

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// UserAddressAttrs 数据库模型属性
type UserAddressAttrs struct {
	Id              int64  `json:"id"`               //收货地址ID
	UserId          int64  `json:"user_id"`          //用户ID
	AdminId         int64  `json:"admin_id"`         //管理员ID
	Name            string `json:"name"`             //收货人名
	Phone           string `json:"phone"`            //电话
	Country         string `json:"country"`          //国家
	ShippingAddress string `json:"shipping_address"` //收货地址
	DoorNumber      int64  `json:"door_number"`      //门牌号
	ZipCode         int64  `json:"zip_code"`         //邮编
	IsDefault       int64  `json:"is_default"`       //是否默认：1是，-1否
	Time            int64  `json:"time"`             //时间
	Status          int64  `json:"status"`           //状态：1在使用，-1已删除
}

const (
	UserAddressDefault        = 1  //	默认地址
	UserAddressNoDefault      = -1 //	普通地址
	UserAddressStatusDelete   = -1 //	已删除
	UserAddressStatusActivate = 1  //	在使用
)

// UserAddress 数据库模型
type UserAddress struct {
	define.Db
}

// NewUserAddress 创建数据库模型
func NewUserAddress(tx *sql.Tx) *UserAddress {
	return &UserAddress{
		database.DbPool.NewDb(tx).Table("user_address"),
	}
}

// AndWhere where条件
func (c *UserAddress) AndWhere(str string, arg ...any) *UserAddress {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *UserAddress) FindOne() *UserAddressAttrs {
	attrs := new(UserAddressAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.UserId, &attrs.AdminId, &attrs.Name, &attrs.Phone, &attrs.Country, &attrs.ShippingAddress, &attrs.DoorNumber, &attrs.ZipCode, &attrs.IsDefault, &attrs.Time, &attrs.Status)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *UserAddress) FindMany() []*UserAddressAttrs {
	data := make([]*UserAddressAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(UserAddressAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.UserId, &tmp.AdminId, &tmp.Name, &tmp.Phone, &tmp.Country, &tmp.ShippingAddress, &tmp.DoorNumber, &tmp.ZipCode, &tmp.IsDefault, &tmp.Time, &tmp.Status)
		data = append(data, tmp)
	})
	return data
}
