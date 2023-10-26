package deposit

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
	Id        int64  `json:"id" validate:"required"`
	Data      string `json:"data"`
	CreatedAt string `json:"created_at"`
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
	model := models.NewUserWalletOrder(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	nowTime := time.Now()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("data=?", params.Data).
		DateTime("created_at=?", params.CreatedAt, location).
		Int64("updated_at=?", nowTime.Unix())

	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id = ?", params.Id).AndWhere("type=?", models.WalletOrderTypeDeposit).Update()
	if err != nil {
		panic(err)
	}
	body.SuccessJSON(w, "ok")
}
