package comment

import (
	"basic/models"
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
	Id         int64  `json:"id" validate:"required"` //评论ID
	StarRating string `json:"star_rating"`            //星级
	Username   string `json:"username"`               //用户名
	UserAvatar string `json:"user_avatar"`            //用户头像
	Comment    string `json:"comment"`                //评论
	Time       int64  `json:"time"`                   //时间
	Status     int64  `json:"status"`                 //状态：1新增，2追加，3回复
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
	model := models.NewUserFavorites(nil)

	// 开启事务
	tx := database.DbPool.GetTx()
	defer tx.Rollback()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("star_rating=?", params.StarRating).
		String("username=?", params.Username).
		String("user_avatar=?", params.UserAvatar).
		String("comment=?", params.Comment).
		Int64("time=?", params.Time).
		Int64("status=?", params.Status)
	// 是否是超级管理员，不是没有权限修改
	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id = ?", params.Id).Update()
	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
