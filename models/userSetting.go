package models

import (
	"database/sql"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// UserSettingAttrs 数据库模型属性
type UserSettingAttrs struct {
	Id      int64  `json:"id"`       //主键
	AdminId int64  `json:"admin_id"` //管理ID
	UserId  int64  `json:"user_id"`  //用户ID
	GroupId int64  `json:"group_id"` //组ID
	Name    string `json:"name"`     //名称
	Type    string `json:"type"`     //类型
	Field   string `json:"field"`    //健名
	Value   string `json:"value"`    //健值
	Data    string `json:"data"`     //数据
}

// UserSetting 数据库模型
type UserSetting struct {
	define.Db
}

// NewUserSetting 创建数据库模型
func NewUserSetting(tx *sql.Tx) *UserSetting {
	return &UserSetting{
		database.DbPool.NewDb(tx).Table("user_setting"),
	}
}

// AndWhere where条件
func (c *UserSetting) AndWhere(str string, arg ...any) *UserSetting {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *UserSetting) FindOne() *UserSettingAttrs {
	attrs := new(UserSettingAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.UserId, &attrs.GroupId, &attrs.Name, &attrs.Type, &attrs.Field, &attrs.Value, &attrs.Data)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *UserSetting) FindMany() []*UserSettingAttrs {
	data := make([]*UserSettingAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(UserSettingAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.UserId, &tmp.GroupId, &tmp.Name, &tmp.Type, &tmp.Field, &tmp.Value, &tmp.Data)
		data = append(data, tmp)
	})
	return data
}

// GetUserSettingValue 获取用户配置值
func (c *UserSetting) GetUserSettingValue(userId int64, field string) string {
	c.AndWhere("user_id=?", userId).AndWhere("field=?", field)

	data := ""
	c.QueryRow(func(row *sql.Row) {
		_ = row.Scan(&data)
	})
	return data
}
