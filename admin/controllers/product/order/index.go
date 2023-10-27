package order

import (
	"basic/models"
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"net/http"
	"strings"
	"time"
)

type indexParams struct {
	ProductName string                 `json:"product_name"`
	AdminName   string                 `json:"admin_name"`
	UserName    string                 `json:"user_name"`
	OrderSn     string                 `json:"order_sn"`
	Status      int64                  `json:"status"`
	Type        int64                  `json:"type"`
	ExpiredAt   *define.RangeTimeParam `json:"expired_at"`
	UpdatedAt   *define.RangeTimeParam `json:"updated_at"`
	CreatedAt   *define.RangeTimeParam `json:"created_at"`
	Pagination  *define.Pagination     `json:"pagination"`
}

type indexData struct {
	models.ProductOrderAttrs
	Data        string `json:"data"`         //	订单数据
	Images      string `json:"images"`       //	图片
	AdminName   string `json:"admin_name"`   //管理员名称
	ProductName string `json:"product_name"` //产品名称
	UserName    string `json:"user_name"`    //用户名称
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取子级包括自己ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	//  实例化模型
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))
	model := models.NewProductOrder(nil)
	model.Db.AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("status>?", models.ProductStatusDelete)

	define.NewFilterEmpty(model.Db).
		String("order_sn like ?", "%"+params.OrderSn+"%").
		Int64("status=?", params.Status).
		Int64("type=?", params.Type).
		RangeTime("expired_at between ? and ?", params.ExpiredAt, location).
		RangeTime("updated_at between ? and ?", params.UpdatedAt, location).
		RangeTime("created_at between ? and ?", params.CreatedAt, location).
		Pagination(params.Pagination)

	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}
	// 用户名称
	if params.UserName != "" {
		model.Db.AndWhere("user_id in (" + strings.Join(models.NewUser(nil).FindUserLikeNameIds(params.UserName), ",") + ")")
	}

	if params.ProductName != "" {
		productIds := models.NewProduct(nil).Field("id").AndWhere("name like ?", "%"+params.ProductName+"%").ColumnString()
		if len(productIds) > 0 {
			model.Db.AndWhere("product_id in (" + strings.Join(productIds, ",") + ")")
		}
	}

	data := make([]*indexData, 0)
	rds := cache.RedisPool.Get()
	defer rds.Close()
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.UserId, &tmp.ProductId, &tmp.OrderSn, &tmp.Money, &tmp.Nums, &tmp.Type, &tmp.Status, &tmp.Data, &tmp.ExpiredAt, &tmp.UpdatedAt, &tmp.CreatedAt)

		// 用户名
		adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}

		// 当前用户信息
		userInfo := models.NewUser(nil).AndWhere("id=?", tmp.UserId).FindOne()
		if userInfo != nil {
			tmp.UserName = userInfo.UserName
		}

		productInfo := models.NewProduct(nil).AndWhere("id=?", tmp.ProductId).FindOne()
		if productInfo != nil {
			tmp.ProductName = locales.Manager.GetAdminLocales(rds, settingAdminId, "zh-CN", productInfo.Name)
			tmp.Images = productInfo.Images
		}
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
