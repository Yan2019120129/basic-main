package models

import (
	"database/sql"
	"encoding/json"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/utils"
)

// AdminMenu 模型
type AdminMenu struct {
	define.Db
}

const (
	AdminMenuStatusActivate = 10 //	激活
	AdminMenuStatusDisabled = -1 //	禁用
)

// AdminMenuAttrs 菜单属性
type AdminMenuAttrs struct {
	Id     int64  `json:"id"`     //	菜单ID
	Name   string `json:"name"`   //	菜单名称
	Parent int64  `json:"parent"` //	父级ID
	Route  string `json:"route"`  //	对应路由
	Sort   int    `json:"sort"`   //	排序
	Status int    `json:"status"` //	状态
	Data   string `json:"data"`   //	菜单数据
}

// AdminMenuAttrsData 菜单属性Data
type AdminMenuAttrsData struct {
	Icon string `json:"icon"` //	菜单图标
	Temp string `json:"temp"` //	对应模版
}

// AdminMenuList 管理员菜单列表
type AdminMenuList struct {
	Name     string              `json:"name"`     //	名称
	Route    string              `json:"route"`    //	路由
	Children []*AdminMenuList    `json:"children"` //	子集
	Data     *AdminMenuAttrsData `json:"data"`     //	数据
}

// NewAdminMenu 创建模型
func NewAdminMenu(tx *sql.Tx) *AdminMenu {
	return &AdminMenu{
		database.DbPool.NewDb(tx).Table("admin_menu"),
	}
}

// AndWhere where条件
func (c *AdminMenu) AndWhere(str string, arg ...any) *AdminMenu {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 单条查询
func (c *AdminMenu) FindOne() *AdminMenuAttrs {
	attrs := new(AdminMenuAttrs)
	c.QueryRow(func(row *sql.Row) {
		menuData := ""
		_ = row.Scan(&attrs.Id, &attrs.Name, &attrs.Parent, &attrs.Route, &attrs.Sort, &attrs.Status, &menuData)
		_ = json.Unmarshal([]byte(menuData), &attrs.Data)
	})
	return attrs
}

// FindMany 查询多条数据
func (c *AdminMenu) FindMany() []*AdminMenuAttrs {
	var data []*AdminMenuAttrs
	c.Query(func(rows *sql.Rows) {
		tmp := new(AdminMenuAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.Name, &tmp.Parent, &tmp.Route, &tmp.Sort, &tmp.Status, &tmp.Data)
		data = append(data, tmp)
	})
	return data
}

// GetAdminMenuList 获取管理员菜单
func (c *AdminMenu) GetAdminMenuList(adminId int64) []*AdminMenuList {
	roleList := NewAdminAuthAssignment(nil).GetAdminRoleList(adminId)
	if len(roleList) > 0 {
		routerList := utils.GetMapValues(NewAdminAuthChild(nil).GetRolesRouteList(roleList))
		if len(routerList) > 0 {
			c.AndWhere("status=?", AdminMenuStatusActivate).OrderBy("parent asc", "sort asc")
			menuList := c.FindMany()
			return c.children(0, menuList, routerList)
		}
	}

	return []*AdminMenuList{}
}

func (c *AdminMenu) children(menuId int64, menuList []*AdminMenuAttrs, routerList []string) []*AdminMenuList {
	var data []*AdminMenuList
	for _, menu := range menuList {
		if menu.Parent == menuId {
			if menu.Route == "" || utils.SliceStringIndexOf("*", routerList) > -1 || utils.SliceStringIndexOf(menu.Route, routerList) > -1 {
				childrenList := c.children(menu.Id, menuList, routerList)
				if menu.Route != "" || len(childrenList) > 0 {
					dataTmp := new(AdminMenuAttrsData)
					_ = json.Unmarshal([]byte(menu.Data), &dataTmp)
					data = append(data, &AdminMenuList{
						Name:     menu.Name,
						Route:    menu.Route,
						Children: childrenList,
						Data:     dataTmp,
					})
				}
			}
		}
	}
	return data
}
