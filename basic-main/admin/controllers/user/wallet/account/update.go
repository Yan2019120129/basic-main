package account

import (
	"basic/models"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type updateParams struct {
	Id         int64  `json:"id" validate:"required"`
	Name       string `json:"name"`
	RealName   string `json:"real_name"`
	CardNumber string `json:"card_number"`
	Address    string `json:"address"`
	Status     int64  `json:"status"`
	Sort       int64  `json:"sort"`
	Data       string `json:"data"`
	CreatedAt  string `json:"created_at"`
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

	//  实例化模型
	model := models.NewUserWalletAccount(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))
	nowTime := time.Now()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("name=?", params.Name).
		String("real_name=?", params.RealName).
		String("card_number=?", params.CardNumber).
		String("address=?", params.Address).
		Int64("status=?", params.Status).
		Int64("sort=?", params.Sort).
		String("data=?", params.Data).
		DateTime("created_at=?", params.CreatedAt, location).
		Int64("updated_at=?", nowTime.Unix())

	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id = ?", params.Id).Update()
	if err != nil {
		panic(err)
	}
	body.SuccessJSON(w, "ok")
}
