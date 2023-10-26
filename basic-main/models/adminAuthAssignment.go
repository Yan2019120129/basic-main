package models

import (
	"database/sql"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// AdminAuthAssignment 模型
type AdminAuthAssignment struct {
	define.Db
}

// AdminAuthAssignmentsAttrs 属性
type AdminAuthAssignmentsAttrs struct {
	ItemName  string `json:"item_name"`  //	名称
	UserId    int64  `json:"user_id"`    //	用户ID
	CreatedAt int64  `json:"created_at"` //	创建时间
}

// NewAdminAuthAssignment 创建模型
func NewAdminAuthAssignment(tx *sql.Tx) *AdminAuthAssignment {
	return &AdminAuthAssignment{
		database.DbPool.NewDb(tx).Table("admin_auth_assignment"),
	}
}

// AndWhere 条件
func (c *AdminAuthAssignment) AndWhere(str string, arg ...any) *AdminAuthAssignment {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单条数据
func (c *AdminAuthAssignment) FindOne() *AdminAuthAssignmentsAttrs {
	attrs := new(AdminAuthAssignmentsAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.ItemName, &attrs.UserId, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// GetAdminRoleList 获取管理员角色列表
func (c *AdminAuthAssignment) GetAdminRoleList(adminId int64) []string {
	var roleList []string
	c.Field("item_name").
		AndWhere("user_id = ?", adminId).
		Query(func(rows *sql.Rows) {
			var role string
			_ = rows.Scan(&role)
			roleList = append(roleList, role)
		})
	return roleList
}
