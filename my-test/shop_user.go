package my_test

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// ShopUserAttrs 表示商店用户信息的结构体
type ShopUserAttrs struct {
	ID           int64   `json:"id"`            // 用户唯一ID
	HeadPortrait string  `json:"head_portrait"` // 用户头像
	Name         string  `json:"name"`          // 用户姓名
	Sex          int     `json:"sex"`           // 用户性别：-1未知，1男，2女
	Birthday     int     `json:"birthday"`      // 用户生日（存储时间戳）
	Email        string  `json:"email"`         // 用户邮箱
	Phone        string  `json:"phone"`         // 用户电话
	UpdateTime   int     `json:"update_time"`   // 用户信息最新修改时间（时间戳格式）
	CreateTime   int     `json:"create_time"`   // 用户注册时间(时间戳格式)
	Status       int     `json:"status"`        // 用户状态：-2删除，-1禁用，10启用
	Balance      float64 `json:"balance"`       // 用户余额
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

// AndWhere 添加 WHERE 子句
func (s *ShopUser) AndWhere(str string, args ...interface{}) *ShopUser {
	s.Db.AndWhere(str, args...)
	return s
}

// FindOne 查询单个商店用户信息
func (s *ShopUser) FindOne() *ShopUserAttrs {
	user := new(ShopUserAttrs)
	s.QueryRow(func(row *sql.Row) {
		//err := row.Scan(&user.ID, &user.HeadPortrait, &user.Name, &user.Sex, &user.Birthday, &user.Email, &user.Phone,
		//	&user.UpdateTime, &user.CreateTime, &user.Status, &user.Balance)
		err := row.Scan(GetFieldAddr(user)...)
		if err != nil {
			user = nil
		}
	})
	return user
}

// FindMany 查询多个商店用户信息
func (s *ShopUser) FindMany() []*ShopUserAttrs {
	data := make([]*ShopUserAttrs, 0)
	s.Query(func(rows *sql.Rows) {
		tmp := new(ShopUserAttrs)
		v := GetFieldAddr(tmp)
		//不使用getAddr需要手动获取地址值：
		//_ = rows.Scan(&tmp.ID, &tmp.HeadPortrait, &tmp.Name, &tmp.Sex, &tmp.Birthday, &tmp.Email, &tmp.Phone,
		//	&tmp.UpdateTime, &tmp.CreateTime, &tmp.Status, &tmp.Balance)
		_ = rows.Scan(v...)
		data = append(data, tmp)
	})
	return data
}
