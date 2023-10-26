package user

import (
	"basic/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type verifyInfoData struct {
	*models.UserVerifyAttrs
	CountryId int64  `json:"country_id"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
}

// VerifyInfo 用户验证信息
func VerifyInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := router.TokenManager.GetContextClaims(r)

	userInfo := models.NewUser(nil).AndWhere("id=?", claims.UserId).FindOne()
	verifyInfo := models.NewUserVerify(nil).AndWhere("user_id=?", userInfo.Id).FindOne()

	data := &verifyInfoData{
		verifyInfo,
		userInfo.CountryId,
		userInfo.Telephone,
		userInfo.Email,
	}
	body.SuccessJSON(w, data)
}
