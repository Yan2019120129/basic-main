package user

import (
	"basic/models"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type updateParams struct {
	Id          int64  `json:"id" validate:"required"`
	ParentId    int64  `json:"parent_id" validate:"omitempty,gt=0"`
	CountryId   int64  `json:"country_id" validate:"omitempty,gt=0"`
	Email       string `json:"email"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Telephone   string `json:"telephone"`
	Sex         int64  `json:"sex" validate:"omitempty,oneof=-1 1 2"`
	Birthday    string `json:"birthday"`
	Password    string `json:"password"`
	SecurityKey string `json:"security_key"`
	Status      int64  `json:"status" validate:"omitempty,oneof=-2 -1 1 10"`
}

// Update 当前用户更新
func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updateParams)
	_ = body.ReadJSON(r, params)

	//  参数验证
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	userInfo := models.NewUser(nil).AndWhere("id=?", params.Id).FindOne()
	if userInfo == nil {
		body.ErrorJSON(w, "用户不存在", -1)
		return
	}

	if params.ParentId > 0 {
		parentInfo := models.NewUser(nil).AndWhere("id=?", params.ParentId).FindOne()
		if parentInfo == nil {
			body.ErrorJSON(w, "邀请人不存在", -1)
			return
		}
	}

	if params.CountryId > 0 {
		countryInfo := models.NewCountry(nil).AndWhere("id=?", params.CountryId).FindOne()
		if countryInfo == nil {
			body.ErrorJSON(w, "国家ID不存在", -1)
			return
		}
	}

	// 	如果状态禁用，那么删除Token
	if params.Status == models.UserStatusDisabled {
		rds := cache.RedisPool.Get()
		defer rds.Close()

		router.TokenManager.DelTokenValue(rds, userInfo.AdminId, userInfo.Id)
	}

	//  实例化模型
	model := models.NewUser(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))
	nowTime := time.Now()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		Int64("parent_id=?", params.ParentId).
		Int64("country_id=?", params.CountryId).
		String("email=?", params.Email).
		String("nickname=?", params.Nickname).
		String("avatar=?", params.Avatar).
		String("telephone=?", params.Telephone).
		Int64("sex=?", params.Sex).
		DateTime("birthday=?", params.Birthday, location).
		String("password=?", utils.PasswordEncrypt(params.Password)).
		String("security_key=?", utils.PasswordEncrypt(params.SecurityKey)).
		Int64("status=?", params.Status).
		Int64("updated_at=?", nowTime.Unix())

	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}

	//  模型增加where条件并更新
	_, err = model.AndWhere("id=?", params.Id).Update()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
