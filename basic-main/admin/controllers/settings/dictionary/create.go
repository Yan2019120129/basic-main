package dictionary

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type createParams struct {
	Alias string `json:"alias" validate:"required"`
	Type  int64  `json:"type" validate:"required,oneof=1 2 10"`
	Name  string `json:"name" validate:"required,max=50"`
	Field string `json:"field" validate:"required,max=50"`
	Value string `json:"value" validate:"required"`
	Data  string `json:"data" validate:"max=255"`
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	//  模型插入数据
	nowTime := time.Now()
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	oldLangInfo := models.NewLang(nil).AndWhere("admin_id=?", adminId).AndWhere("alias=?", params.Alias).AndWhere("status>?", models.LangStatusDisabled).FindOne()
	if oldLangInfo == nil {
		body.ErrorJSON(w, "语言别名不存在", -1)
		return
	}

	_, err = models.NewLangDictionary(nil).
		Field("admin_id", "type", "alias", "name", "field", "value", "data", "created_at").
		Args(adminId, params.Type, params.Alias, params.Name, params.Field, params.Value, params.Data, nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	// 如果更新对类型是接口，或者数据翻译， 那么重载语言配置
	if params.Value != "" && params.Field != "" {
		rds := cache.RedisPool.Get()
		defer rds.Close()
		locales.Manager.SetAdminLocales(rds, adminId, params.Alias, params.Field, params.Value)
	}
	body.SuccessJSON(w, "ok")
}
