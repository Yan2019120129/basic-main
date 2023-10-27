package role

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type updateParams struct {
	Name        string          `json:"name" validate:"required"`
	NewName     string          `json:"new_name"`
	Description string          `json:"description"`
	Authority   map[string]bool `json:"authority"`
}

// Update 更新角色
func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updateParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	//  验证角色是否存在
	authItemRoleInfo := models.NewAdminAuthItem(nil).AndWhere("name=?", params.Name).AndWhere("type=?", models.AdminAuthItemTypeManage).FindOne()
	if authItemRoleInfo == nil {
		body.ErrorJSON(w, "角色不存在", -1)
		return
	}

	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	//	修改角色信息
	model := models.NewAdminAuthItem(tx)
	nowTime := time.Now()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("name=?", params.NewName).
		String("description=?", params.Description).
		Int64("updated_at=?", nowTime.Unix())
	_, err = model.AndWhere("name=?", params.Name).AndWhere("type=?", models.AdminAuthItemTypeManage).Update()
	if err != nil {
		panic(err)
	}

	// 更新角色权限
	if len(params.Authority) > 0 {
		// 删除原有的权限
		_, err = models.NewAdminAuthChild(nil).AndWhere("parent=?", authItemRoleInfo.Name).AndWhere("type=?", models.AdminAuthItemTypeRouteName).Delete()
		if err != nil {
			panic(err)
		}
		for auth, check := range params.Authority {
			if check {
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
