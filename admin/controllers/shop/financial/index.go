package financial

import (
	"basic/models"
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"net/http"
	"strings"
)

type indexParams struct {
	AdminName      string             `json:"admin_name"`      // 管理员名
	ProductName    string             `json:"product_name"`    //商品名
	Profit         float64            `json:"profit"`          //利润
	UnitPrice      float64            `json:"unit_price"`      //单价
	WholesalePrice float64            `json:"wholesale_price"` //批发价
	Pagination     *define.Pagination `json:"pagination"`      //	分页
}

type indexData struct {
	models.FinancialStatisticsAttrs
	AdminName   string `json:"admin_name"`
	ProductName string `json:"product_name"` //商品名
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取子级包括自己ID
	// 从token 中获取管理员ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	// 获取管理员ID，用于判断该管理员是否有权限获取信息
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)

	model := models.NewFinancialStatistics(nil)
	// 判断用户收藏中的管理员ID在管理员中
	model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	// 查询条件
	define.NewFilterEmpty(model.Db).
		// 添加product_name，date字段的条件
		String("product_name=?", params.ProductName).
		Float64("profit=?", params.Profit).
		Float64("unit_price=?", params.UnitPrice).
		Float64("wholesale_price=?", params.WholesalePrice).
		Pagination(params.Pagination)

	// 管理员名称ƒ
	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	// 获取数据
	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.ProductId, &tmp.Profit, &tmp.UnitPrice, &tmp.WholesalePrice)
		// 获取对应的管理员名
		adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}

		// 获取对应的用户名
		commodityInfo := models.NewCommodity(nil).AndWhere("id=?", tmp.ProductId).FindOne()
		if commodityInfo != nil {
			tmp.ProductName = commodityInfo.Name
		}
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
