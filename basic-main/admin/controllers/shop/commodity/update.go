package commodity

import (
	"basic/models"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
	"net/http"
	"strings"
)

type updateParams struct {
	Id                      int64   `json:"id" validate:"required,gt=0"` //类目ID
	ProductImage            string  `json:"product_image"`               //商品图片
	Name                    string  `json:"name"`                        //名称
	PurchasePrice           float64 `json:"purchase_price"`              //进货价
	SellingPrice            float64 `json:"selling_price"`               //出售价
	Stock                   int64   `json:"stock"`                       //库存
	SalesVolume             int64   `json:"sales_volume"`                //出货量
	Status                  int64   `json:"status"`                      //状态 -1删除 -2下架 1在售
	Operation               int64   `json:"operation"`                   //操作
	SpecificationAttributes string  `json:"specification_attributes"`    //规格属性
	Description             string  `json:"description"`                 //描述
	Brand                   string  `json:"brand"`                       //品牌
}

func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updateParams)
	_ = body.ReadJSON(r, params)
	//  参数验证
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 从token获取管理员ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	model := models.NewCommodity(nil)

	// 开启事务
	tx := database.DbPool.GetTx()
	defer tx.Rollback()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("product_image=?", params.ProductImage).
		String("name=?", params.Name).
		Float64("purchase_price=?", params.PurchasePrice).
		Float64("selling_price=?", params.SellingPrice).
		Int64("stock=?", params.Stock).
		Int64("sales_volume=?", params.SalesVolume).
		Int64("status=?", params.Status).
		Int64("operation=?", params.Operation).
		String("specification_attributes=?", params.SpecificationAttributes).
		String("description=?", params.Description).
		String("brand=?", params.Brand)
	// 是否是超级管理员，不是没有权限修改
	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id = ?", params.Id).Update()
	fmt.Println("err", err)
	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
