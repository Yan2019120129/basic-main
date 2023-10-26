package models

import (
	"database/sql"

	"github.com/gomodule/redigo/redis"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/locales"
)

const (
	ProductCategoryStatusActivate = 10 //	分类状态启用
	ProductCategoryStatusDisabled = -1 //	分类状态禁用
	ProductCategoryStatusDelete   = -2 //	分类状态删除
	ProductCategoryRecommend      = 10 //	分类推荐
	ProductCategoryTypeDefault    = 1  //	分类类型默认
)

var ProductCategoryTypeList map[int64]string = map[int64]string{
	ProductCategoryTypeDefault: "默认类型",
}

type ProductCateogryTypesAttrs struct {
	Label string `json:"label"`
	Value int64  `json:"value"`
}

// ProductCategoryAttrs 数据库模型属性
type ProductCategoryAttrs struct {
	Id        int64  `json:"id"`         //主键
	ParentId  int64  `json:"parent_id"`  //分类上级ID
	AdminId   int64  `json:"admin_id"`   //管理员ID
	Type      int64  `json:"type"`       //类型 1收益商品 2电商商品 10区块链
	Name      string `json:"name"`       //标题
	Image     string `json:"image"`      //封面
	Sort      int64  `json:"sort"`       //排序
	Status    int64  `json:"status"`     //状态 -2删除 -1禁用 10启用
	Recommend int64  `json:"recommend"`  //推荐 -1关闭 10推荐
	Data      string `json:"data"`       //数据
	UpdatedAt int64  `json:"updated_at"` //更新时间
	CreatedAt int64  `json:"created_at"` //创建时间
}

// ProductCategory 数据库模型
type ProductCategory struct {
	define.Db
}

// NewProductCategory 创建数据库模型
func NewProductCategory(tx *sql.Tx) *ProductCategory {
	return &ProductCategory{
		database.DbPool.NewDb(tx).Table("product_category"),
	}
}

// AndWhere where条件
func (c *ProductCategory) AndWhere(str string, arg ...any) *ProductCategory {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *ProductCategory) FindOne() *ProductCategoryAttrs {
	attrs := new(ProductCategoryAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.ParentId, &attrs.AdminId, &attrs.Type, &attrs.Name, &attrs.Image, &attrs.Sort, &attrs.Status, &attrs.Recommend, &attrs.Data, &attrs.UpdatedAt, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *ProductCategory) FindMany() []*ProductCategoryAttrs {
	data := make([]*ProductCategoryAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(ProductCategoryAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.ParentId, &tmp.AdminId, &tmp.Type, &tmp.Name, &tmp.Image, &tmp.Sort, &tmp.Status, &tmp.Recommend, &tmp.Data, &tmp.UpdatedAt, &tmp.CreatedAt)
		data = append(data, tmp)
	})
	return data
}

// InheritAdminProduct 继承管理员产品
func (c *ProductCategory) InheritAdminProduct(adminId, newAdminId, adminParentId, newAdminParentId int64) error {
	productCategoryModel := NewProductCategory(nil)
	productCategoryModel.AndWhere("admin_id=?", adminId).AndWhere("parent_id=?", adminParentId)
	productCategoryList := productCategoryModel.FindMany()

	for i := 0; i < len(productCategoryList); i++ {
		categoryItem := productCategoryList[i]

		productModel := NewProduct(nil)
		productModel.AndWhere("category_id=?", categoryItem.Id)
		productList := productModel.FindMany()

		// 插入管理员分类
		newCategoryId, err := NewProductCategory(c.GetTx()).Field("parent_id", "admin_id", "type", "name", "image", "sort", "status", "recommend", "data", "updated_at", "created_at").
			Args(newAdminParentId, newAdminId, categoryItem.Type, categoryItem.Name, categoryItem.Image, categoryItem.Sort, categoryItem.Status, categoryItem.Recommend, categoryItem.Data, categoryItem.UpdatedAt, categoryItem.CreatedAt).
			Insert()
		if err != nil {
			return err
		}

		for j := 0; j < len(productList); j++ {
			product := productList[j]

			//	获取对应的资产ID
			adminAssetsInfo := NewAssets(c.GetTx()).AndWhere("id=?", product.AssetsId).FindOne()

			var currentAssetsId int64 = 0
			if adminAssetsInfo != nil {
				//	当前管理的资产ID
				newAdminAssetsInfo := NewAssets(c.GetTx()).AndWhere("admin_id=?", newAdminId).AndWhere("name=?", adminAssetsInfo.Name).FindOne()
				currentAssetsId = newAdminAssetsInfo.Id
			}

			// 插入 分类产品
			_, err = NewProduct(c.GetTx()).Field("admin_id", "category_id", "assets_id", "name", "images", "money", "type", "sort", "status", "recommend", "sales", "nums", "used", "total", "data", "describes", "updated_at", "created_at").
				Args(newAdminId, newCategoryId, currentAssetsId, product.Name, product.Images, product.Money, product.Type, product.Sort, product.Status, product.Recommend, product.Sales, product.Nums, product.Used, product.Total, product.Data, product.Describes, product.UpdatedAt, product.CreatedAt).
				Insert()
			if err != nil {
				return err
			}
		}

		//	往下执行
		c.InheritAdminProduct(adminId, newAdminId, categoryItem.Id, newCategoryId)
	}
	return nil
}

// ProductCategoryTypes 产品分类列表
func ProductCategoryTypes(rds redis.Conn, parentId, settingAdminId int64, sep string) []*ProductCateogryTypesAttrs {
	data := make([]*ProductCateogryTypesAttrs, 0)
	NewProductCategory(nil).Field("id", "parent_id", "name", "type").AndWhere("parent_id=?", parentId).Query(func(rows *sql.Rows) {
		var parentId, categoryType int64
		attrs := new(ProductCateogryTypesAttrs)
		err := rows.Scan(&attrs.Value, &parentId, &attrs.Label, &categoryType)
		if err == nil {
			attrs.Label = sep + locales.Manager.GetAdminLocales(rds, settingAdminId, "zh-CN", attrs.Label) + "[" + ProductCategoryTypeList[categoryType] + "]"
			data = append(data, attrs)

			if parentId > 0 {
				data = append(data, ProductCategoryTypes(rds, parentId, settingAdminId, "---")...)
			}
		}
	})

	return data
}
