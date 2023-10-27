package verify

import (
	"basic/models"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type indexParams struct {
	AdminName  string                 `json:"admin_name"`
	UserName   string                 `json:"username"`
	Status     int64                  `json:"status"`
	Type       int64                  `json:"type"`
	RealName   string                 `json:"real_name"`
	IdNumber   string                 `json:"id_number"`
	Address    string                 `json:"address"`
	Data       string                 `json:"data"`
	DateTime   *define.RangeTimeParam `json:"updated_at"`
	Pagination *define.Pagination     `json:"pagination"`
}

type indexData struct {
	models.UserVerifyAttrs
	AdminName string `json:"admin_name"` //管理员名称
	UserName  string `json:"username"`   //用户名
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取子级包括自己ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)

	model := models.NewUserVerify(nil)
	model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	define.NewFilterEmpty(model.Db).
		Int64("status=?", params.Status).
		Int64("type=?", params.Type).
		String("real_name=?", params.RealName).
		String("id_number=?", params.IdNumber).
		String("address=?", params.Address).
		String("data=?", params.Data).
		RangeTime("created_at between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)

	// 管理员名称
	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	// 用户名称
	if params.UserName != "" {
		model.Db.AndWhere("user_id in (" + strings.Join(models.NewUser(nil).FindUserLikeNameIds(params.UserName), ",") + ")")
	}

	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.UserId, &tmp.Type, &tmp.RealName, &tmp.IdNumber, &tmp.IdPhoto1, &tmp.IdPhoto2, &tmp.IdPhoto3, &tmp.Address, &tmp.Data, &tmp.Status, &tmp.CreatedAt, &tmp.UpdatedAt)
		// 当前用户信息
		userInfo := models.NewUser(nil).AndWhere("id=?", tmp.UserId).FindOne()
		if userInfo != nil {
			tmp.UserName = userInfo.UserName
		}

		// 当前管理员信息
		adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
