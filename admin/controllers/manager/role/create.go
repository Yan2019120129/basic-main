package role

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type createParams struct {
	Name        string          `json:"name" validate:"required"`
	Authority   map[string]bool `json:"authority"`
	Description string          `json:"description"`
}

// Create 创建角色
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	nowTime := time.Now()
	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	_, err = models.NewAdminAuthItem(tx).
		Field("name", "type", "description", "created_at", "updated_at").
		Args(params.Name, models.AdminAuthItemTypeManage, params.Description, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	// 新增角色权限
	for auth, check := range params.Authority {
		if check {
			authInfo := models.NewAdminAuthItem(nil).AndWhere("name=?", auth).AndWhere("type=?", models.AdminAuthItemTypeRouteName).FindOne()
			if authInfo != nil && check {
				_, err = models.NewAdminAuthChild(tx).Field("parent", "child", "type").
					Args(params.Name, auth, models.AdminAuthItemTypeRouteName).
					Insert()
				if err != nil {
					panic(err)
				}
			}
		}
	}

	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
