package index

import (
	"basic/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

// Info 管理员信息
func Info(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	adminInfo := models.NewAdminUser(nil).AndWhere("id=?", adminId).FindOne()

	data := &userInfo{
		Id:         adminInfo.Id,
		Username:   adminInfo.UserName,
		Nickname:   adminInfo.Nickname,
		Email:      adminInfo.Email,
		Avatar:     adminInfo.Avatar,
		Money:      adminInfo.Money,
		Data:       adminInfo.Data,
		InviteCode: models.NewUserInvite(nil).GetInviteCode(adminInfo.Id, 0),
		ExpiredAt:  adminInfo.ExpiredAt,
		UpdatedAt:  adminInfo.UpdatedAt,
	}

	body.SuccessJSON(w, data)
}
