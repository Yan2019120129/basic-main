package models

import (
	"database/sql"
	"errors"

	"github.com/ip2location/ip2location-go/v9"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

const (
	CountryStatusActivate = 10
	CountryStatusDisabled = -1
)

// CountryAttrs 数据库模型属性
type CountryAttrs struct {
	Id        int64  `json:"id"`         //主键
	AdminId   int64  `json:"admin_id"`   //管理员ID
	LangId    int64  `json:"lang_id"`    //语言ID
	Name      string `json:"name"`       //名称
	Alias     string `json:"alias"`      //别名
	Iso1      string `json:"iso1"`       //	ISO3166-1
	Icon      string `json:"icon"`       //图标
	Code      string `json:"code"`       //区号
	Sort      int64  `json:"sort"`       //排序
	Status    int64  `json:"status"`     //状态 -1禁用｜10启用
	Data      string `json:"data"`       //数据
	CreatedAt int64  `json:"created_at"` //创建时间
}

// Country 数据库模型
type Country struct {
	define.Db
}

// NewCountry 创建数据库模型
func NewCountry(tx *sql.Tx) *Country {
	return &Country{
		database.DbPool.NewDb(tx).Table("country"),
	}
}

// AndWhere where条件
func (c *Country) AndWhere(str string, arg ...any) *Country {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *Country) FindOne() *CountryAttrs {
	attrs := new(CountryAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.LangId, &attrs.Name, &attrs.Alias, &attrs.Iso1, &attrs.Icon, &attrs.Code, &attrs.Sort, &attrs.Status, &attrs.Data, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *Country) FindMany() []*CountryAttrs {
	var data []*CountryAttrs
	c.Query(func(rows *sql.Rows) {
		tmp := new(CountryAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.LangId, &tmp.Name, &tmp.Alias, &tmp.Iso1, &tmp.Icon, &tmp.Code, &tmp.Sort, &tmp.Status, &tmp.Data, &tmp.CreatedAt)
		data = append(data, tmp)
	})
	return data
}

// AutoUserCountry 自动获取用户国家
func (c *Country) AutoUserCountry(adminId int64, userIP string) (int64, error) {
	if userIP == "0.0.0.0" {
		return 0, nil
	}

	ip2locationDB, err := ip2location.OpenDB(LiteDB1Path)
	if err != nil {
		return 0, err
	}
	defer ip2locationDB.Close()
	results, err := ip2locationDB.Get_all(userIP)
	if err != nil {
		return 0, err
	}
	if results.Country_short == "-" {
		return 0, errors.New("ip2location no info")
	}

	c.AndWhere("admin_id=?", adminId).AndWhere("iso1=?", results.Country_short)
	countryInfo := c.FindOne()
	if countryInfo == nil {
		return 0, errors.New("country no info")
	}
	return countryInfo.Id, nil
}
