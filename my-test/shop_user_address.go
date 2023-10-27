package my_test

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// UserAddressAttrs 表示用户收货地址的结构体
type UserAddressAttrs struct {
	ID              int    `json:"id"`               // 收货地址ID
	UserID          int    `json:"user_id"`          // 用户ID
	Name            string `json:"name"`             // 收货人名
	Phone           string `json:"phone"`            // 电话
	Country         string `json:"country"`          // 国家
	ShippingAddress string `json:"shipping_address"` // 收货地址
	DoorNumber      int    `json:"door_number"`      // 门牌号
	ZipCode         int    `json:"zip_code"`         // 邮编
	IsDefault       int    `json:"is_default"`       // 是否默认
	Time            int    `json:"time"`             // 时间
	Status          int    `json:"status"`           // 状态
}

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

// AndWhere 添加 WHERE 子句
func (u *UserAddress) AndWhere(str string, args ...interface{}) *UserAddress {
	u.Db.AndWhere(str, args...)
	return u
}

// FindOne 查询单个用户收货地址信息
func (u *UserAddress) FindOne() *UserAddressAttrs {
	address := new(UserAddressAttrs)
	u.QueryRow(func(row *sql.Row) {
		err := row.Scan(&address.ID, &address.UserID, &address.Name, &address.Phone, &address.Country, &address.ShippingAddress,
			&address.DoorNumber, &address.ZipCode, &address.IsDefault, &address.Time, &address.Status)
		if err != nil {
			address = nil
		}
	})
	return address
}

// FindMany 查询多个用户收货地址信息
func (u *UserAddress) FindMany() []*UserAddressAttrs {
	data := make([]*UserAddressAttrs, 0)
	u.Query(func(rows *sql.Rows) {
		tmp := new(UserAddressAttrs)
		_ = rows.Scan(&tmp.ID, &tmp.UserID, &tmp.Name, &tmp.Phone, &tmp.Country, &tmp.ShippingAddress,
			&tmp.DoorNumber, &tmp.ZipCode, &tmp.IsDefault, &tmp.Time, &tmp.Status)
		data = append(data, tmp)
	})
	return data
}
