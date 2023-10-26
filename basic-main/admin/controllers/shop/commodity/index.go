package commodity

import (
	"basic/models"
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"net/http"
	"strings"
	"time"
)

type indexParams struct {
	AdminName               string                 `json:"admin_name"`               //管理员名
	StoreName               string                 `json:"store_name"`               //店铺名
	CategoryName            string                 `json:"category_name"`            //类目ID
	CommodityName           string                 `json:"class_name"`               //分类ID
	ProductImage            string                 `json:"product_image"`            //商品图片
	Name                    string                 `json:"name"`                     //名称
	PurchasePrice           float64                `json:"purchase_price"`           //进货价
	SellingPrice            float64                `json:"selling_price"`            //出售价
	Stock                   string                 `json:"stock"`                    //库存
	SalesVolume             string                 `json:"sales_volume"`             //出货量
	Status                  string                 `json:"status"`                   //状态
	Operation               string                 `json:"operation"`                //操作
	SpecificationAttributes string                 `json:"specification_attributes"` //规格属性
	Description             string                 `json:"description"`              //描述
	Brand                   string                 `json:"brand"`                    //品牌
	State                   string                 `json:"state"`                    //状态
	Time                    *define.RangeTimeParam `json:"time"`                     // 时间戳
	Pagination              *define.Pagination     `json:"pagination"`               //	分页
}

type indexData struct {
	models.CommodityAttrs
	AdminName     string `json:"admin_name"`
	StoreName     string `json:"store_name"`    //店铺名
	CategoryName  string `json:"category_name"` //类目ID
	CommodityName string `json:"commodity_id"`  //分类ID
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取子级包括自己ID
	// 从token 中获取管理员ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	// 获取管理员ID，用于判断该管理员是否有权限获取信息
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(adminId, "site_timezone"))

	model := models.NewCommodity(nil)
	// 判断用户收藏中的管理员ID在管理员中
	model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	// 查询条件
	define.NewFilterEmpty(model.Db).
		// 添加product_name，date字段的条件
		String("product_image=?", params.ProductImage).
		String("name=?", params.Name).
		Float64("purchase_price=?", params.PurchasePrice).
		Float64("selling_price=?", params.SellingPrice).
		String("stock=?", params.Stock).
		String("sales_volume=?", params.SalesVolume).
		String("status=?", params.Status).
		String("operation=?", params.Operation).
		String("specification_attributes=?", params.SpecificationAttributes).
		String("description=?", params.Description).
		String("brand=?", params.Brand).
		String("state=?", params.State).
		RangeTime("time between ? and ?", params.Time, location).
		Pagination(params.Pagination)

	// 管理员名称
	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}
	// 获取数据
	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.StoreId, &tmp.AdminId, &tmp.ProductImage, &tmp.Name, &tmp.PurchasePrice, &tmp.SellingPrice, &tmp.Stock, &tmp.SalesVolume, &tmp.Status, &tmp.Operation, &tmp.CategoryId, &tmp.CommodityId, &tmp.SpecificationAttributes, &tmp.Description, &tmp.Brand, &tmp.Time)
		// 获取对应的管理员名
		adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}
		//店铺名
		storeInfo := models.NewStore(nil).AndWhere("id=?", tmp.StoreId).FindOne()
		if storeInfo != nil {
			tmp.StoreName = storeInfo.StoreName
		}
		//类目ID
		categoryInfo := models.NewCategory(nil).AndWhere("id=?", tmp.CategoryId).FindOne()
		if adminInfo != nil {
			tmp.CategoryName = categoryInfo.Name
		}
		//分类ID
		commodityIdInfo := models.NewCommodity(nil).AndWhere("id=?", tmp.CommodityId).FindOne()
		if adminInfo != nil {
			tmp.CategoryName = commodityIdInfo.Name
		}
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
