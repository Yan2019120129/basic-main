package manager

import (
	"basic/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type searchParams struct {
	ParentName string                 `json:"parent_name"`
	Username   string                 `json:"username"`
	Nickname   string                 `json:"nickname"`
	Email      string                 `json:"email"`
	Status     int64                  `json:"status"`
	Pagination *define.Pagination     `json:"pagination"`
	DateTime   *define.RangeTimeParam `json:"updated_at"`
}

type adminAttrs struct {
	ParentName string              `json:"parent_name"` //	父级名称
	InviteCode string              `json:"invite_code"` //	邀请码
	Role       string              `json:"role"`        //角色
	Roles      map[string]bool     `json:"roles"`       //角色列表
	Data       *router.TokenParams `json:"data"`        //Token参数
	Nums       int64               `json:"nums"`        //	代理人数
	Template   string              `json:"template"`    //	模版名称
	models.AdminUserAttrs
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(searchParams)
	_ = body.ReadJSON(r, params)

	//  实例化模型
	model := models.NewAdminUser(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))
	model.AndWhere("id in ("+strings.Join(adminIds, ",")+")").AndWhere("status>?", models.AdminUserStatusDelete)

	//  模型增加where条件并分页
	define.NewFilterEmpty(model.Db).
		String("username like ?", "%"+params.Username+"%").
		String("nickname like ?", "%"+params.Nickname+"%").
		String("email like ?", "%"+params.Email+"%").
		Int64("status = ?", params.Status).
		RangeTime("created_at between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)

	if params.ParentName != "" {
		model.Db.AndWhere("parent_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.ParentName), ",") + ")")
	}

	//  查询数据
	data := make([]*adminAttrs, 0)
	model.Field("id", "parent_id", "username", "email", "nickname", "avatar", "money", "status", "domain", "data", "expired_at", "updated_at", "created_at").
		Query(func(rows *sql.Rows) {
			tmp := new(adminAttrs)
			var adminDataStr, adminExtraStr string
			_ = rows.Scan(&tmp.Id, &tmp.ParentId, &tmp.UserName, &tmp.Email, &tmp.Nickname,
				&tmp.Avatar, &tmp.Money, &tmp.Status, &tmp.Domain, &adminDataStr, &adminExtraStr, &tmp.ExpiredAt, &tmp.UpdatedAt, &tmp.CreatedAt)
			_ = json.Unmarshal([]byte(adminDataStr), &tmp.Data)

			adminExtra := new(models.AdminUserExtraAttrs)
			_ = json.Unmarshal([]byte(adminExtraStr), &adminExtra)
			tmp.Template = adminExtra.Template
			tmp.Nums = adminExtra.Nums

			adminTmpInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.ParentId).FindOne()
			if adminTmpInfo != nil {
				tmp.ParentName = adminTmpInfo.UserName
			}
			// 邀请码
			tmp.InviteCode = models.NewUserInvite(nil).GetInviteCode(tmp.Id, 0)
			// 获取角色
			adminRoles := models.NewAdminAuthAssignment(nil).GetAdminRoleList(tmp.Id)
			tmp.Role = strings.Join(adminRoles, ",")
			tmp.Roles = models.NewAdminAuthItem(nil).GetAdminRoleCheckedList(adminId, adminRoles)
			data = append(data, tmp)
		})

	//  统计数量
	count := model.Count()

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: count,
	})
}
