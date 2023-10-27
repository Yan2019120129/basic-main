package dictionary

import (
	"basic/models"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type updateParams struct {
	Id    int64  `json:"id" validate:"required"`
	Alias string `json:"alias"`
	Type  int64  `json:"type" validate:"omitempty,oneof=1 2 10"`
	Name  string `json:"name" validate:"max=50"`
	Field string `json:"field" validate:"max=50"`
	Value string `json:"value"`
	Data  string `json:"data" validate:"max=255"`
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
	model := models.NewLangDictionary(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("alias=?", params.Alias).
		Int64("type=?", params.Type).
		String("name=?", params.Name).
		String("field=?", params.Field).
		String("data=?", params.Data)

	if params.Value != "" {
		if params.Value == " " {
			params.Value = strings.Replace(params.Value, " ", "", -1)
		}
		model.Value("value=?").Args(params.Value)
	}

	//  模型增加where条件并更新
	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id=?", params.Id).Update()
	if err != nil {
		panic(err)
	}

	// 获取当前信息
	dictionaryInfo := models.NewLangDictionary(nil).AndWhere("id=?", params.Id).FindOne()
	if dictionaryInfo != nil {
		rds := cache.RedisPool.Get()
		defer rds.Close()
		locales.Manager.SetAdminLocales(rds, adminId, dictionaryInfo.Alias, dictionaryInfo.Field, dictionaryInfo.Value)
	}

	body.SuccessJSON(w, "ok")
}
