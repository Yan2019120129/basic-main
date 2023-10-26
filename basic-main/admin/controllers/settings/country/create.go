package country

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type createParams struct {
	LangId int64  `json:"lang_id"`
	Name   string `json:"name" validate:"required,max=50"`
	Alias  string `json:"alias" validate:"required,max=50"`
	Iso1   string `json:"iso1" validate:"required,max=50"`
	Icon   string `json:"icon" validate:"required,max=255"`
	Code   string `json:"code" validate:"required,max=50"`
	Data   string `json:"data"`
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)
	//  验证参数
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	//  模型插入数据
	nowTime := time.Now()
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	//	不能重复的别名
	oldCountryInfo := models.NewCountry(nil).AndWhere("admin_id=?", adminId).AndWhere("alias=?", params.Alias).AndWhere("status>=?", models.CountryStatusDisabled).FindOne()
	if oldCountryInfo != nil {
		body.ErrorJSON(w, "当前国家别名已存在", -1)
		return
	}

	_, err = models.NewCountry(nil).
		Field("admin_id", "lang_id", "name", "alias", "iso1", "icon", "code", "data", "created_at").
		Args(adminId, params.LangId, params.Name, params.Alias, params.Iso1, params.Icon, params.Code, params.Data, nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}
	body.SuccessJSON(w, "ok")
}
