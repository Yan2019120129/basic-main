package payment

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
	AdminName   string                 `json:"admin_name"`
	Name        string                 `json:"name"`
	AccountName string                 `json:"account_name"`
	AccountCode string                 `json:"account_code"`
	Type        int64                  `json:"type"`
	Mode        int64                  `json:"mode"`
	Status      int64                  `json:"status"`
	UpdatedAt   *define.RangeTimeParam `json:"updated_at"`
	Pagination  *define.Pagination     `json:"pagination"`
}

type indexItems struct {
	AdminName string `json:"admin_name"` //管理名称
	models.WalletPaymentAttrs
}

// Index 角色列表
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取请求管理员的uid
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	model := models.NewWalletPayment(nil)
	model.Db.AndWhere("status<>?", models.WalletPaymentStatusDelete)
	define.NewFilterEmpty(model.Db).
		String("name like ?", "%"+params.Name+"%").
		String("account_name like ?", "%"+params.AccountName+"%").
		String("account_code like ?", "%"+params.AccountCode+"%").
		Int64("status=?", params.Status).
		Int64("type=?", params.Type).
		Int64("mode=?", params.Mode).
		RangeTime("created_at between ? and ?", params.UpdatedAt, location).
		Pagination(params.Pagination)

	model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	data := make([]*indexItems, 0)
	model.Field("id", "admin_id", "icon", "mode", "type", "name", "account_name", "account_code", "sort", "status", "description", "data", "expand", "created_at", "updated_at").
		Query(func(rows *sql.Rows) {
			tmp := new(indexItems)
			_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.Icon, &tmp.Mode, &tmp.Type,
				&tmp.Name, &tmp.AccountName, &tmp.AccountCode, &tmp.Sort, &tmp.Status, &tmp.Description, &tmp.Data, &tmp.Expand, &tmp.CreatedAt, &tmp.UpdatedAt)
			adminModel := models.NewAdminUser(nil)
			adminModel.AndWhere("id=?", tmp.AdminId)
			adminInfoTmp := adminModel.FindOne()
			if adminInfoTmp != nil {
				tmp.AdminName = adminInfoTmp.UserName
			}
			data = append(data, tmp)
		})

	count := model.Count()
	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: count,
	})
}
