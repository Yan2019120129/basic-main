package manager

import (
	"basic/models"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type updateExtraParams struct {
	Id       int64  `json:"id"`       //	管理ID
	Template string `json:"template"` //	模版名称
	Nums     int64  `json:"nums"`     //	管理人数
}

// UpdateExtra 更新扩展字段
func UpdateExtra(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updateExtraParams)
	_ = body.ReadJSON(r, params)

	//  参数验证
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 判断管理员是否存在
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	updateAdminInfo := models.NewAdminUser(nil).AndWhere("id=?", params.Id).FindOne()
	if updateAdminInfo == nil {
		body.ErrorJSON(w, "更新管理员不存在", -1)
		return
	}

	//	如果不是超级管理员不能修改
	if adminId != models.AdminUserSupermanId {
		body.ErrorJSON(w, "不是超级管理员不能修改", -1)
		return
	}

	adminExtra := new(models.AdminUserExtraAttrs)
	adminExtra.Template = params.Template
	adminExtra.Nums = params.Nums
	adminExtraBytes, _ := json.Marshal(adminExtra)

	model := models.NewAdminUser(nil)
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("extra=?", string(adminExtraBytes))

	_, err = model.AndWhere("id=?", params.Id).Update()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
