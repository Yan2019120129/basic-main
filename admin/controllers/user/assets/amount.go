package assets

import (
	"basic/models"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type amountParams struct {
	UserName string  `json:"username" validate:"required"`
	Type     int64   `json:"type" validate:"required,oneof=101 102"`
	AssetsId int64   `json:"assets_id" validate:"required"`
	Money    float64 `json:"money" validate:"required,number"`
}

// Amount 修改用户余额
func Amount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(amountParams)
	_ = body.ReadJSON(r, params)
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	adminId := router.TokenManager.GetContextClaims(r).AdminId
	userMode := models.NewUser(nil)
	userMode.AndWhere("username=?", params.UserName)
	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		userMode.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	userInfo := userMode.FindOne()
	if userInfo == nil {
		body.ErrorJSON(w, "用户不存在", -1)
		return
	}

	// 资产是否存在
	assetsInfo := models.NewAssets(nil).AndWhere("id=?", params.AssetsId).FindOne()
	if assetsInfo == nil {
		body.ErrorJSON(w, "资产不存在", -1)
		return
	}

	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	switch params.Type {
	case models.UserBillTypeAssetsSystemDeposit:
		models.UserAssetsDeposit(tx, userInfo, assetsInfo, 0, params.Type, params.Money)
	case models.UserBillTypeAssetsSystemDeduction:
		models.UserAssetsSpend(tx, userInfo, assetsInfo, 0, params.Type, params.Money)
	}

	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
