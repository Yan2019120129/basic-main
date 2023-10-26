package statistical

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
	AdminName             string                 `json:"admin_name"`              //管理员名
	VisitorCount          int64                  `json:"visitor_count"`           //访客数
	OrderCount            int64                  `json:"order_count"`             //订单数
	Earnings              float64                `json:"earnings"`                //收益
	ShopFavoritesCount    int64                  `json:"shop_favorites_count"`    //店铺收藏量
	Credit                int64                  `json:"credit"`                  //信用
	ProductFavoritesCount int64                  `json:"product_favorites_count"` //商品收藏量
	ProductCount          int64                  `json:"product_count"`           //商品数量
	Time                  *define.RangeTimeParam `json:"time"`                    // 时间戳
	Pagination            *define.Pagination     `json:"pagination"`              //	分页
}

type indexData struct {
	models.ShopStatisticalRecordAttrs
	AdminName string `json:"admin_name"`
	UserName  string `json:"username"`
	ShopName  int64  `json:"shop_name"` //店铺ID

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

	model := models.NewShopStatisticalRecord(nil)
	// 判断用户收藏中的管理员ID在管理员中
	model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	// 查询条件
	define.NewFilterEmpty(model.Db).
		// 添加product_name，date字段的条件
		Int64("visitor_count=?", params.VisitorCount).
		Int64("order_count=?", params.OrderCount).
		Float64("earnings=?", params.Earnings).
		Int64("shop_favorites_count=?", params.ShopFavoritesCount).
		Int64("credit=?", params.Credit).
		Int64("product_favorites_count=?", params.ProductFavoritesCount).
		Int64("product_count=?", params.ProductCount).
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
		_ = rows.Scan(&tmp.Id, &tmp.ShopId, &tmp.VisitorCount, &tmp.OrderCount, &tmp.Earnings, &tmp.ShopFavoritesCount, &tmp.Credit, &tmp.ProductFavoritesCount, &tmp.ProductCount, &tmp.Time)
		// 获取对应的管理员名
		adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}
		// 获取对应的商品名
		storeInfo := models.NewStore(nil).AndWhere("id=?", tmp.ShopId).FindOne()
		if storeInfo != nil {
			tmp.UserName = storeInfo.StoreName
		}
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
