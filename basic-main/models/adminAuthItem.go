package models

import (
	"database/sql"
	"fmt"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/utils"
)

type AdminAuthItemAttrs struct {
	Name        string         `json:"name"`        //名称
	Type        int64          `json:"type"`        //类型
	Description string         `json:"description"` //描述
	RuleName    string         `json:"rule_name"`   //规则名称
	Data        sql.NullString `json:"data"`        //数据
	CreatedAt   int64          `json:"created_at"`  //创建时间
	UpdatedAt   int64          `json:"updated_at"`  //更新时间
}

const (
	// AdminAuthItemTypeManage 管理员名称
	AdminAuthItemTypeManage = 1
	// AdminAuthItemTypeRoute 请求路由
	AdminAuthItemTypeRoute = 2
	// AdminAuthItemTypeRouteName 请求路由名称
	AdminAuthItemTypeRouteName = 3
)

type AdminAuthItem struct {
	define.Db
}

func NewAdminAuthItem(tx *sql.Tx) *AdminAuthItem {
	return &AdminAuthItem{
		database.DbPool.NewDb(tx).Table("admin_auth_item"),
	}
}

// AndWhere where条件
func (c *AdminAuthItem) AndWhere(str string, arg ...any) *AdminAuthItem {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑数据
func (c *AdminAuthItem) FindOne() *AdminAuthItemAttrs {
	attrs := new(AdminAuthItemAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Name, &attrs.Type, &attrs.Description, &attrs.RuleName, &attrs.Data, &attrs.CreatedAt, &attrs.UpdatedAt)
		if err != nil {
			fmt.Println(err)
			attrs = nil
		}
	})
	return attrs
}

// FindMany 多个多条数据
func (c *AdminAuthItem) FindMany() []*AdminAuthItemAttrs {
	data := make([]*AdminAuthItemAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		attrs := new(AdminAuthItemAttrs)
		_ = rows.Scan(&attrs.Name, &attrs.Type, &attrs.Description, &attrs.RuleName, &attrs.Data, &attrs.CreatedAt, &attrs.UpdatedAt)
		data = append(data, attrs)
	})
	return data
}

// GetAdminRoleCheckedList 获取管理员角色选中列表
func (c *AdminAuthItem) GetAdminRoleCheckedList(adminId int64, adminRoles []string) map[string]bool {
	adminInfo := NewAdminUser(nil).AndWhere("id=?", adminId).FindOne()

	c.Db.AndWhere("type = ?", AdminAuthItemTypeManage)
	// 如果自己是超级管理员
	if adminInfo.Id == AdminUserSupermanId {
		c.Db.AndWhere("name <> ?", "超级管理员").AndWhere("name <> ?", "默认组长管理员")
	}

	if adminInfo.ParentId == AdminUserSupermanId {
		currentRoleInfo := NewAdminAuthAssignment(nil).AndWhere("user_id=?", adminInfo.Id).FindOne()
		if currentRoleInfo != nil {
			c.Db.AndWhere("name like ?", string([]rune(currentRoleInfo.ItemName)[:2])+"%")
		}
		c.Db.AndWhere("name <> ?", "超级管理员").AndWhere("name <> ?", "默认后台管理员")
	}

	data := map[string]bool{}
	c.Field("name").Query(func(rows *sql.Rows) {
		var name string
		_ = rows.Scan(&name)
		data[name] = false
		if utils.SliceStringIndexOf(name, adminRoles) > -1 {
			data[name] = true
		}
	})
	return data
}
