package models

import (
	"database/sql"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

const (
	ProductNumsUnlimited  = -1 //	没有限制购买
	ProductStatusActivate = 10 //	状态激活
	ProductStatusDisabled = -1 //	状态禁用
	ProductStatusDelete   = -2 //	状态删除
	ProductRecommend      = 10 //	产品推荐
)

// ProductAttrs 数据库模型属性
type ProductAttrs struct {
	Id         int64   `json:"id"`          //主键
	AdminId    int64   `json:"admin_id"`    //管理员ID
	CategoryId int64   `json:"category_id"` //类目ID
	AssetsId   int64   `json:"assets_id"`   //资产ID
	Name       string  `json:"name"`        //标题
	Images     string  `json:"images"`      //图片列表
	Money      float64 `json:"money"`       //金额
	Type       int64   `json:"type"`        //类型 1默认
	Sort       int64   `json:"sort"`        //排序
	Status     int64   `json:"status"`      //状态 -2删除 -1禁用 10启用
	Recommend  int64   `json:"recommend"`   //推荐 -1关闭 10推荐
	Sales      int64   `json:"sales"`       //销售量
	Nums       int64   `json:"nums"`        //限购 -1无限
	Used       int64   `json:"used"`        //已使用
	Total      int64   `json:"total"`       //总数
	Data       string  `json:"data"`        //数据
	Describes  string  `json:"describes"`   //数据
	UpdatedAt  int64   `json:"updated_at"`  //更新时间
	CreatedAt  int64   `json:"created_at"`  //创建时间
}

// Product 数据库模型
type Product struct {
	define.Db
}

// NewProduct 创建数据库模型
func NewProduct(tx *sql.Tx) *Product {
	return &Product{
		database.DbPool.NewDb(tx).Table("product"),
	}
}

// AndWhere where条件
func (c *Product) AndWhere(str string, arg ...any) *Product {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *Product) FindOne() *ProductAttrs {
	attrs := new(ProductAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.CategoryId, &attrs.AssetsId, &attrs.Name, &attrs.Images, &attrs.Money, &attrs.Type, &attrs.Sort, &attrs.Status, &attrs.Recommend, &attrs.Sales, &attrs.Nums, &attrs.Used, &attrs.Total, &attrs.Data, &attrs.Describes, &attrs.UpdatedAt, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *Product) FindMany() []*ProductAttrs {
	data := make([]*ProductAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(ProductAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.CategoryId, &tmp.AssetsId, &tmp.Name, &tmp.Images, &tmp.Money, &tmp.Type, &tmp.Sort, &tmp.Status, &tmp.Recommend, &tmp.Sales, &tmp.Nums, &tmp.Used, &tmp.Total, &tmp.Data, &tmp.Describes, &tmp.UpdatedAt, &tmp.CreatedAt)
		data = append(data, tmp)
	})
	return data
}
