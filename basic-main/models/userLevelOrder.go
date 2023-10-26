package models

import (
	"database/sql"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

const (
	UserLevelOrderStatusActivate = 10
)

type UserLevelOrderAttrs struct {
	Id        int64  `json:"id"`         //主键
	AdminId   int64  `json:"admin_id"`   //管理员ID
	UserId    int64  `json:"user_id"`    //用户ID
	LevelId   int64  `json:"level_id"`   //等级ID
	Data      string `json:"data"`       //数据
	Status    int64  `json:"status"`     //状态 -2删除 -1禁用 10启用
	CreatedAt int64  `json:"created_at"` //创建时间
	UpdatedAt int64  `json:"updated_at"` //更新时间
}

type UserLevelOrder struct {
	define.Db
}

func NewUserLevelOrder(tx *sql.Tx) *UserLevelOrder {
	return &UserLevelOrder{
		database.DbPool.NewDb(tx).Table("user_level_order"),
	}
}

// AndWhere where条件
func (c *UserLevelOrder) AndWhere(str string, arg ...any) *UserLevelOrder {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 单条信息
func (c *UserLevelOrder) FindOne() *UserLevelOrderAttrs {
	attrs := new(UserLevelOrderAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.UserId, &attrs.LevelId, &attrs.Data, &attrs.Status, &attrs.CreatedAt, &attrs.UpdatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}
