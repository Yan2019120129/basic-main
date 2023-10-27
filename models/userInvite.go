package models

import (
	"database/sql"
	"time"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/utils"
)

type UserInviteAttrs struct {
	Id        int64  `json:"id"`         //主键
	AdminId   int64  `json:"admin_id"`   //管理员ID
	UserId    int64  `json:"user_id"`    //用户ID
	Code      string `json:"code"`       //邀请码
	Status    int64  `json:"status"`     //状态 -2删除 -1禁用 10启用
	Data      string `json:"data"`       //数据
	CreatedAt int64  `json:"created_at"` //创建时间
}

type UserInvite struct {
	define.Db
}

func NewUserInvite(tx *sql.Tx) *UserInvite {
	return &UserInvite{
		database.DbPool.NewDb(tx).Table("user_invite"),
	}
}

// AndWhere where条件
func (c *UserInvite) AndWhere(str string, arg ...any) *UserInvite {
	c.Db.AndWhere(str, arg...)
	return c
}

func (c *UserInvite) FindOne() *UserInviteAttrs {
	attrs := new(UserInviteAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.UserId, &attrs.Code, &attrs.Status, &attrs.Data, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// GetInviteCode 获取邀请码
func (c *UserInvite) GetInviteCode(adminId, userId int64) string {
	c.AndWhere("admin_id=?", adminId).AndWhere("user_id=?", userId)
	inviteInfo := c.FindOne()
	if inviteInfo != nil {
		return inviteInfo.Code
	}

	nowTime := time.Now()
	code := utils.NewRandom().SetNumberRunes().String(6)
	_, err := NewUserInvite(nil).Field("admin_id", "user_id", "code", "created_at").
		Args(adminId, userId, code, nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}
	return code
}
