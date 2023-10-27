package models

import (
	"database/sql"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

const (
	LocalesRedisName                = "_locales"
	LangDictionaryTypeApiTip        = 1  //	接口提示
	LangDictionaryTypeDataTranslate = 2  //	数据翻译
	LangDictionaryTypeHomeTranslate = 10 //	前台翻译

)

// LangDictionaryAttrs 数据库模型属性
type LangDictionaryAttrs struct {
	Id        int64  `json:"id"`         //主键
	AdminId   int64  `json:"admin_id"`   //管理员ID
	Type      int64  `json:"type"`       //类型 1接口提示 2前台翻译 3数据翻译
	Alias     string `json:"alias"`      //别名
	Name      string `json:"name"`       //名称
	Field     string `json:"field"`      //键
	Value     string `json:"value"`      //值
	Data      string `json:"data"`       //数据
	CreatedAt int64  `json:"created_at"` //创建时间
}

// LangDictionary 数据库模型
type LangDictionary struct {
	define.Db
}

// NewLangDictionary 创建数据库模型
func NewLangDictionary(tx *sql.Tx) *LangDictionary {
	return &LangDictionary{
		database.DbPool.NewDb(tx).Table("lang_dictionary"),
	}
}

// AndWhere where条件
func (c *LangDictionary) AndWhere(str string, arg ...any) *LangDictionary {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *LangDictionary) FindOne() *LangDictionaryAttrs {
	attrs := new(LangDictionaryAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.Type, &attrs.Alias, &attrs.Name, &attrs.Field, &attrs.Value, &attrs.Data, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *LangDictionary) FindMany() []*LangDictionaryAttrs {
	data := make([]*LangDictionaryAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(LangDictionaryAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.Type, &tmp.Alias, &tmp.Name, &tmp.Field, &tmp.Value, &tmp.Data, &tmp.CreatedAt)
		data = append(data, tmp)
	})
	return data
}

// GetAdminLocales 获取管理本地语言
func GetAdminLocales() map[int64]map[string]map[string]string {
	data := map[int64]map[string]map[string]string{}

	NewLangDictionary(nil).Field("admin_id", "alias", "field", "value").Query(func(rows *sql.Rows) {
		var tmpAdminId int64
		var tmpAlias, tmpField, tmpValue string
		_ = rows.Scan(&tmpAdminId, &tmpAlias, &tmpField, &tmpValue)

		if _, ok := data[tmpAdminId]; !ok {
			data[tmpAdminId] = map[string]map[string]string{}
		}
		if _, ok := data[tmpAdminId][tmpAlias]; !ok {
			data[tmpAdminId][tmpAlias] = map[string]string{}
		}

		data[tmpAdminId][tmpAlias][tmpField] = tmpValue
	})

	return data
}
