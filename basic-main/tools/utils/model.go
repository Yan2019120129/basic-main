package utils

import (
	"fmt"
	"strings"
)

var ModelTemplate = `package models

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// {{TableNameHump}}Attrs 数据库模型属性
type {{TableNameHump}}Attrs struct {{{TableAttrsList}}}

// {{TableNameHump}} 数据库模型
type {{TableNameHump}} struct {
	define.Db
}

// New{{TableNameHump}} 创建数据库模型
func New{{TableNameHump}}(tx *sql.Tx) *{{TableNameHump}} {
	return &{{TableNameHump}}{
		database.DbPool.NewDb(tx).Table("{{TableName}}"),
	}
}

// AndWhere where条件
func (c *{{TableNameHump}}) AndWhere(str string, arg ...any) *{{TableNameHump}} {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *{{TableNameHump}}) FindOne() *{{TableNameHump}}Attrs {
	attrs := new({{TableNameHump}}Attrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan({{ScanAttrs}})
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *{{TableNameHump}}) FindMany() []*{{TableNameHump}}Attrs {
	data := make([]*{{TableNameHump}}Attrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new({{TableNameHump}}Attrs)
		_ = rows.Scan({{ScanTmpAttrs}})
		data = append(data, tmp)
	})
	return data
}`

// GenerateModel 生成模型
func (c *Operate) GenerateModel(path string, tableName string) {
	c.Columns(tableName)

	tableNameHump := c.HumpName(tableName)
	content := ModelTemplate

	content = strings.ReplaceAll(content, "{{TableName}}", tableName)
	content = strings.ReplaceAll(content, "{{TableNameHump}}", tableNameHump)
	content = strings.ReplaceAll(content, "{{TableAttrsList}}", "\n"+c.ColumnsOpt("attrs", "\t", "\n")+"\n")
	content = strings.ReplaceAll(content, "{{ScanAttrs}}", c.ColumnsOpt("", "&attrs.", ", "))
	content = strings.ReplaceAll(content, "{{ScanTmpAttrs}}", c.ColumnsOpt("", "&tmp.", ", "))

	fileName := strings.ToLower(string(tableNameHump[0])) + tableNameHump[1:]
	err := c.SetFilePath(path).SetFileText(content).SaveFile(fileName + ".go")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("生成数据库模型" + tableName + "成功")
}
