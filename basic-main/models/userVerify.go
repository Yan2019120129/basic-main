package models

import (
	"database/sql"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

const (
	UserVerifyTypeIdCard     = 1
	UserVerifyStatusPending  = 10
	UserVerifyStatusComplete = 20
)

// UserVerifyAttrs 数据库模型属性
type UserVerifyAttrs struct {
	Id        int64  `json:"id"`         //主键
	AdminId   int64  `json:"admin_id"`   //管理员ID
	UserId    int64  `json:"user_id"`    //用户ID
	Type      int64  `json:"type"`       //类型 1身份证 2护照
	RealName  string `json:"real_name"`  //真实姓名
	IdNumber  string `json:"id_number"`  //证件号码
	IdPhoto1  string `json:"id_photo1"`  //证件照1
	IdPhoto2  string `json:"id_photo2"`  //证件照2
	IdPhoto3  string `json:"id_photo3"`  //证件照3
	Address   string `json:"address"`    //地址
	Data      string `json:"data"`       //数据
	Status    int64  `json:"status"`     //状态 -1拒绝｜10审核｜20通过
	CreatedAt int64  `json:"created_at"` //创建时间
	UpdatedAt int64  `json:"updated_at"` //更新时间
}

// UserVerify 数据库模型
type UserVerify struct {
	define.Db
}

// NewUserVerify 创建数据库模型
func NewUserVerify(tx *sql.Tx) *UserVerify {
	return &UserVerify{
		database.DbPool.NewDb(tx).Table("user_verify"),
	}
}

// AndWhere where条件
func (c *UserVerify) AndWhere(str string, arg ...any) *UserVerify {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *UserVerify) FindOne() *UserVerifyAttrs {
	attrs := new(UserVerifyAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.UserId, &attrs.Type, &attrs.RealName, &attrs.IdNumber, &attrs.IdPhoto1, &attrs.IdPhoto2, &attrs.IdPhoto3, &attrs.Address, &attrs.Data, &attrs.Status, &attrs.CreatedAt, &attrs.UpdatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *UserVerify) FindMany() []*UserVerifyAttrs {
	data := make([]*UserVerifyAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(UserVerifyAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.UserId, &tmp.Type, &tmp.RealName, &tmp.IdNumber, &tmp.IdPhoto1, &tmp.IdPhoto2, &tmp.IdPhoto3, &tmp.Address, &tmp.Data, &tmp.Status, &tmp.CreatedAt, &tmp.UpdatedAt)
		data = append(data, tmp)
	})
	return data
}
