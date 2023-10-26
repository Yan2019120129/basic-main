package utils

import (
	"errors"
	"os"
	"strings"

	"github.com/so68/zfeng"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
)

// Columns 表属性
type Columns struct {
	Field      string      `json:"field"`
	Type       string      `json:"type"`
	Collation  interface{} `json:"collation"`
	Null       string      `json:"null"`
	Key        string      `json:"key"`
	Default    interface{} `json:"default"`
	Extra      string      `json:"extra"`
	Privileges string      `json:"privileges"`
	Comment    string      `json:"comment"`
}

// InitTable 初始化表
type InitTable struct {
	Name        string //	表名
	Comment     string //	注释
	CreateTable string //	创建表语句
	InsertTable string //	插入表语句
}

// Operate 数据库操作
type Operate struct {
	TableName  string     //	表名
	FilePath   string     //	文件路径
	FileText   string     //	文件内容
	ColumnList []*Columns //	表字段列表
}

// Permission 权限
type Permission struct {
	Roles            []string            //	角色列表
	RoleFilterRouter map[string][]string //	角色过滤路由
	RoleOnlyRouter   map[string][]string //	角色拥有路由
}

// NewOperate 创建数据库操作
func NewOperate() *Operate {
	return &Operate{
		ColumnList: make([]*Columns, 0),
	}
}

// Reload 刷新数据库
func (c *Operate) Reload(tableList []*InitTable) *Operate {
	databaseSQL := map[string]string{
		"dropDatabase":   "DROP DATABASE IF EXISTS " + zfeng.ConfInfo.Database.Dbname,
		"createDatabase": "CREATE DATABASE " + zfeng.ConfInfo.Database.Dbname,
	}

	databaseName := zfeng.ConfInfo.Database.Dbname
	zfeng.ConfInfo.Database.Dbname = ""
	database.InitializationDb(zfeng.ConfInfo.Database)
	//	删除数据库
	_, err := database.DbPool.GetDb().Exec(databaseSQL["dropDatabase"])
	if err != nil {
		panic(err)
	}

	//	创建数据库
	_, err = database.DbPool.GetDb().Exec(databaseSQL["createDatabase"])
	if err != nil {
		panic(err)
	}

	//	创建所有的表
	zfeng.ConfInfo.Database.Dbname = databaseName
	database.InitializationDb(zfeng.ConfInfo.Database)
	for i := 0; i < len(tableList); i++ {
		_, err = database.DbPool.GetDb().Exec(tableList[i].CreateTable)
		if err != nil {
			panic(err)
		}
		if tableList[i].InsertTable != "" {
			_, err = database.DbPool.GetDb().Exec(tableList[i].InsertTable)
			if err != nil {
				panic(err)
			}
		}
	}
	return c
}

// LoadRouterList 加载权限表
func (c *Operate) LoadRouterList(routerList []*router.Handle, permission *Permission) *Operate {
	itemList := []string{"('超级管理员', 1), ('*', 2), ('所有权限', 3)"}
	childList := []string{"('所有权限', '*', 2), ('超级管理员', '所有权限', 3)"}

	// 	添加角色
	for _, role := range permission.Roles {
		itemList = append(itemList, "('"+role+"', 1)")
	}

	//	分配角色权限
	for _, handle := range routerList {
		if handle.RouteAuth {
			itemList = append(itemList, "('"+handle.Route+"', 2), ('"+handle.Name+"', 3)")
			childList = append(childList, "('"+handle.Name+"', '"+handle.Route+"', 2)")

			for _, role := range permission.Roles {
				// 	如果属于过滤路由, 那么过滤
				if _, ok := permission.RoleFilterRouter[role]; ok {
					if utils.SliceStringIndexOf(handle.Name, permission.RoleFilterRouter[role]) == -1 {
						childList = append(childList, "('"+role+"', '"+handle.Name+"', 3)")
					}
				} else {
					if onlyRouter, ok := permission.RoleOnlyRouter[role]; ok {
						if utils.SliceStringIndexOf(handle.Name, onlyRouter) > -1 {
							childList = append(childList, "('"+role+"', '"+handle.Name+"', 3)")
						}
					}
				}
			}
		}
	}

	//	新增数据 RBAC item表
	insertItem := "INSERT INTO admin_auth_item(name, type) VALUES" + strings.Join(itemList, ",") + ";"
	_, err := database.DbPool.GetDb().Exec(insertItem)
	if err != nil {
		panic(err)
	}
	//	新增数据	RBAC child表
	insertChild := "INSERT INTO admin_auth_child(parent, child, type) VALUES" + strings.Join(childList, ",") + ";"
	_, err = database.DbPool.GetDb().Exec(insertChild)
	if err != nil {
		panic(err)
	}
	insertAssignment := "INSERT INTO admin_auth_assignment(item_name, user_id) VALUES ('超级管理员', 1);"
	_, err = database.DbPool.GetDb().Exec(insertAssignment)
	if err != nil {
		panic(err)
	}
	return c
}

// Columns 表字段列
func (c *Operate) Columns(name string) []*Columns {
	c.ColumnList = []*Columns{}
	rows, err := database.DbPool.GetDb().Query("show full columns from `" + name + "`")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		tmp := new(Columns)
		err = rows.Scan(&tmp.Field, &tmp.Type, &tmp.Collation, &tmp.Null, &tmp.Key, &tmp.Default, &tmp.Extra, &tmp.Privileges, &tmp.Comment)
		if err != nil {
			panic(err)
		}
		c.ColumnList = append(c.ColumnList, tmp)
	}
	return c.ColumnList
}

// ColumnsOpt 表字段操作
func (c *Operate) ColumnsOpt(opt string, firstStr string, joinStr string) string {
	var strList []string
	switch opt {
	case "attrs":
		for i := 0; i < len(c.ColumnList); i++ {
			item := c.ColumnList[i]
			strList = append(strList, firstStr+c.HumpName(item.Field)+"\t"+c.conversionGoType(item.Type)+"\t`json:\""+item.Field+"\"`"+"\t//"+item.Comment)
		}
		return strings.Join(strList, joinStr)
	default:
		for i := 0; i < len(c.ColumnList); i++ {
			item := c.ColumnList[i]
			strList = append(strList, firstStr+c.HumpName(item.Field))
		}
		return strings.Join(strList, joinStr)
	}
}

// SaveFile 保存文件
func (c *Operate) SaveFile(fileName string) error {
	dir, _ := os.Getwd()

	if !utils.PathExists(dir + c.FilePath) {
		utils.PathMkdirAll(dir + c.FilePath)
	}

	filePath := dir + c.FilePath + "/" + fileName
	if utils.PathExists(filePath) {
		return errors.New("【" + fileName + "】已存在...")
	}

	err := utils.FileWrite(filePath, []byte(c.FileText))
	if err != nil {
		panic(err)
	}
	return nil
}

// SetFilePath 设置文件路径
func (c *Operate) SetFilePath(path string) *Operate {
	c.FilePath = path
	return c
}

// SetFileText 设置文件内容
func (c *Operate) SetFileText(content string) *Operate {
	c.FileText = content
	return c
}

// HumpName 驼峰名称
func (c *Operate) HumpName(name string) string {
	var str string
	nameList := strings.Split(name, "_")

	for i := 0; i < len(nameList); i++ {
		str += strings.ToUpper(nameList[i][:1]) + nameList[i][1:]
	}
	return str
}

func (c *Operate) conversionGoType(dbType string) string {
	switch {
	case strings.Index(dbType, "int") > -1:
		return "int64"
	case strings.Index(dbType, "decimal") > -1:
		return "float64"
	case strings.Index(dbType, "tinyint") > -1, strings.Index(dbType, "smallint") > -1:
		return "int"
	}
	return "string"
}
