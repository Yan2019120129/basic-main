package verify

import (
	"basic/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type createParams struct {
	UserName string `json:"username" validate:"required"`
	Type     int64  `json:"type" validate:"required,oneof=1 2"`
	RealName string `json:"real_name" validate:"required"`
	IdNumber string `json:"id_number" validate:"required"`
	IdPhoto1 string `json:"id_photo1" validate:"required"`
	IdPhoto2 string `json:"id_photo2" validate:"required"`
	IdPhoto3 string `json:"id_photo3" validate:"required"`
	Address  string `json:"address" validate:"required"`
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	userInfo := models.NewUser(nil).AndWhere("username=?", params.UserName).FindOne()
	if userInfo == nil {
		body.ErrorJSON(w, "用户名不存在", -1)
		return
	}

	// 新增验证方法
	nowTime := time.Now()
	_, err = models.NewUserVerify(nil).
		Field("admin_id", "user_id", "type", "real_name", "id_number", "id_photo1", "id_photo2", "id_photo3", "address", "created_at", "updated_at").
		Args(userInfo.AdminId, userInfo.Id, params.Type, params.RealName, params.IdNumber, params.IdPhoto1, params.IdPhoto2, params.IdPhoto3, params.Address, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
