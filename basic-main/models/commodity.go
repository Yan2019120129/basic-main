package models

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// CommodityAttrs 数据库模型属性
type CommodityAttrs struct {
	Id                      int64   `json:"id"`                       //商品ID
	AdminId                 int64   `json:"admin_id"`                 //管理员ID
	StoreId                 int64   `json:"store_id"`                 //店铺ID
	ProductImage            string  `json:"product_image"`            //商品图片
	Name                    string  `json:"name"`                     //名称
	PurchasePrice           float64 `json:"purchase_price"`           //进货价
	SellingPrice            float64 `json:"selling_price"`            //出售价
	Stock                   int64   `json:"stock"`                    //库存
	SalesVolume             int64   `json:"sales_volume"`             //出货量
	Status                  int64   `json:"status"`                   //状态 -1删除 -2下架 1在售
	Operation               int64   `json:"operation"`                //操作
	CategoryId              int64   `json:"category_id"`              //类目ID
	CommodityId             int64   `json:"commodity_id"`             //分类ID
	SpecificationAttributes string  `json:"specification_attributes"` //规格属性
	Description             string  `json:"description"`              //描述
	Brand                   string  `json:"brand"`                    //品牌
	Time                    int64   `json:"time"`                     //时间
}

const (
	CommodityStatusDelete        = -1 //-1删除
	CommodityStatusUndercarriage = -2 //2下架
	CommodityStatusOnSale        = 1  //2下架
)

// Commodity 数据库模型
type Commodity struct {
	define.Db
}

// NewCommodity 创建数据库模型
func NewCommodity(tx *sql.Tx) *Commodity {
	return &Commodity{
		database.DbPool.NewDb(tx).Table("commodity"),
	}
}

// AndWhere where条件
func (c *Commodity) AndWhere(str string, arg ...any) *Commodity {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *Commodity) FindOne() *CommodityAttrs {
	attrs := new(CommodityAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.StoreId, &attrs.AdminId, &attrs.ProductImage, &attrs.Name, &attrs.PurchasePrice, &attrs.SellingPrice, &attrs.Stock, &attrs.SalesVolume, &attrs.Status, &attrs.Operation, &attrs.CategoryId, &attrs.CommodityId, &attrs.SpecificationAttributes, &attrs.Description, &attrs.Brand, &attrs.Time)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *Commodity) FindMany() []*CommodityAttrs {
	data := make([]*CommodityAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(CommodityAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.StoreId, &tmp.AdminId, &tmp.ProductImage, &tmp.Name, &tmp.PurchasePrice, &tmp.SellingPrice, &tmp.Stock, &tmp.SalesVolume, &tmp.Status, &tmp.Operation, &tmp.CategoryId, &tmp.CommodityId, &tmp.SpecificationAttributes, &tmp.Description, &tmp.Brand, &tmp.Time)
		data = append(data, tmp)
	})
	return data
}
