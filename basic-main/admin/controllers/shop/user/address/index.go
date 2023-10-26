package address

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
	UserName        string                 `json:"username"`         //用户名
	AdminName       string                 `json:"admin_name"`       //管理员名
	Country         string                 `json:"country"`          //国家
	ShippingAddress string                 `json:"shipping_address"` //收货地址
	IsDefault       int64                  `json:"is_default"`       //是否默认
	Status          int64                  `json:"status"`           //状态
	DateTime        *define.RangeTimeParam `json:"created_at"`       // 时间戳
	Pagination      *define.Pagination     `json:"pagination"`       //	分页
}

type indexData struct {
	models.UserAddressAttrs
	AdminName string `json:"admin_name"`
	UserName  string `json:"username"`
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

	model := models.NewUserAddress(nil)
	// 判断用户收藏中的管理员ID在管理员中
	model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	// 查询条件
	define.NewFilterEmpty(model.Db).
		// 添加product_name，date字段的条件
		String("country=?", params.Country).
		String("shipping_address=?", params.ShippingAddress).
		Int64("is_default=?", params.IsDefault).
		Int64("status=?", params.Status).
		RangeTime("date between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)

	// 管理员名称
	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}
	// 用户名称
	if params.UserName != "" {
		model.Db.AndWhere("user_id in (" + strings.Join(models.NewUser(nil).FindUserLikeNameIds(params.UserName), ",") + ")")
	}

	// 获取数据
	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.UserId, &tmp.AdminId, &tmp.Name, &tmp.Phone, &tmp.Country, &tmp.ShippingAddress, &tmp.DoorNumber, &tmp.ZipCode, &tmp.IsDefault, &tmp.Time, &tmp.Status)
		// 获取对应的管理员名
		adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}

		// 获取对应的用户名
		userInfo := models.NewUser(nil).AndWhere("id=?", tmp.UserId).FindOne()
		if userInfo != nil {
			tmp.UserName = userInfo.UserName
		}
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
