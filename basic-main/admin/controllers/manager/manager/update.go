package manager

import (
	"basic/models"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type updateParams struct {
	Id          int64               `json:"id" validate:"required"`
	Roles       map[string]bool     `json:"roles"`
	Domain      string              `json:"domain"`
	Email       string              `json:"email"`
	Nickname    string              `json:"nickname"`
	Avatar      string              `json:"avatar"`
	Password    string              `json:"password"`
	SecurityKey string              `json:"security_key"`
	Data        *router.TokenParams `json:"data"`
	Status      int64               `json:"status" validate:"omitempty,oneof=-1 10"`
	ExpiredAt   string              `json:"expired_at"`
}

// Update 当前管理员更新
func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updateParams)
	_ = body.ReadJSON(r, params)

	//  参数验证
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 判断管理员是否存在
	updateAdminInfo := models.NewAdminUser(nil).AndWhere("id=?", params.Id).FindOne()
	if updateAdminInfo == nil {
		body.ErrorJSON(w, "更新管理员不存在", -1)
		return
	}

	//	判断域名是否存在
	domainInfo := models.NewAdminUser(nil).AndWhere("domain like ?", "%"+params.Domain+"%").AndWhere("admin_id<>?", updateAdminInfo.Id).FindOne()
	if domainInfo != nil {
		body.ErrorJSON(w, "域名已存在", -1)
		return
	}

	//  实例化模型
	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	model := models.NewAdminUser(tx)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	nowTime := time.Now()
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	//	如果不是超级管理员不能修改
	if adminId != models.AdminUserSupermanId {
		body.ErrorJSON(w, "不是超级管理员不能修改", -1)
		return
	}

	//  模型设置更新   过滤参数
	adminDataStr := ""
	if params.Data != nil {
		adminDataBytes, _ := json.Marshal(params.Data)
		adminDataStr = string(adminDataBytes)

		//	更新Token缓存
		if adminId == models.AdminUserSupermanId {
			rds := cache.RedisPool.Get()
			defer rds.Close()

			router.TokenManager.SetTokenParams(rds, models.TokenParamsPrefix(models.AdminPrefixTokenKey, updateAdminInfo.Id), params.Data)
		}
	}
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("email=?", params.Email).
		String("domain=?", params.Domain).
		String("nickname=?", params.Nickname).
		String("avatar=?", params.Avatar).
		Int64("status=?", params.Status).
		String("password=?", utils.PasswordEncrypt(params.Password)).
		String("security_key=?", utils.PasswordEncrypt(params.SecurityKey)).
		String("data=?", adminDataStr).
		DateTime("expired_at=?", params.ExpiredAt, location).
		Int64("updated_at=?", nowTime.Unix())

	//  模型增加where条件并更新
	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("parent_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id = ?", updateAdminInfo.Id).Update()
	if err != nil {
		panic(err)
	}

	// 	验证角色是否存在，并且添加角色
	if adminId == models.AdminUserSupermanId {
		adminRoles := models.NewAdminAuthAssignment(nil).GetAdminRoleList(adminId)
		adminRolesList := models.NewAdminAuthItem(nil).GetAdminRoleCheckedList(adminId, adminRoles)
		if len(params.Roles) > 0 {
			_, err = models.NewAdminAuthAssignment(tx).AndWhere("user_id=?", updateAdminInfo.Id).Delete()
			if err != nil {
				panic(err)
			}
		}
		for role, checked := range params.Roles {
			if _, ok := adminRolesList[role]; ok && checked {
				_, err = models.NewAdminAuthAssignment(tx).Field("item_name", "user_id", "created_at").
					Args(role, updateAdminInfo.Id, nowTime.Unix()).
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
