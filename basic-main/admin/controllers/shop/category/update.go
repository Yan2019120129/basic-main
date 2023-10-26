package shop_category

import (
	"basic/models"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
	"net/http"
	"strings"
)

type updateParams struct {
	Id      int64  `json:"id" validate:"required,gt=0"` //类目ID
	Image   string `json:"image"`                       //类目图片
	Name    string `json:"name"`                        //种类名称
	Date    int64  `json:"date"`                        //时间
	Status  int64  `json:"status"`                      //状态 -2删除 -1禁用 10启用
	Operate int64  `json:"operate"`                     //操作
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

	// 从token获取管理员ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	model := models.NewStore(nil)

	// 开启事务
	tx := database.DbPool.GetTx()
	defer tx.Rollback()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("image=?", params.Image).
		String("name=?", params.Name).
		Int64("date=?", params.Date).
		Int64("status=?", params.Status).
		Int64("operate=?", params.Operate)
	// 是否是超级管理员，不是没有权限修改
	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id = ?", params.Id).Update()
	fmt.Println("err", err)
	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
