package user

import (
	"basic/models"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type indexParams struct {
	AdminName   string                 `json:"admin_name"`
	ParentName  string                 `json:"parent_name"`
	CountryName string                 `json:"country_name"`
	Username    string                 `json:"username"`
	Email       string                 `json:"email"`
	Telephone   string                 `json:"telephone"`
	Nickname    string                 `json:"nickname"`
	Sex         int64                  `json:"sex"`
	Birthday    *define.RangeTimeParam `json:"birthday"`
	Type        int64                  `json:"type"`
	Status      int64                  `json:"status"`
	DateTime    *define.RangeTimeParam `json:"updated_at"`
	Pagination  *define.Pagination     `json:"pagination"`
}

type indexData struct {
	AdminName   string `json:"admin_name"`   //管理员名称
	ParentName  string `json:"parent_name"`  //上级名称
	CountryName string `json:"country_name"` //国家名称
	InviteCode  string `json:"invite_code"`  //邀请码
	models.UserAttrs
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取子级包括自己ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	//  实例化模型
	model := models.NewUser(nil)
	model.AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("status<>?", models.UserStatusDelete).AndWhere("type <> ?", models.UserTypeTemporary)

	//  模型增加where条件并分页
	define.NewFilterEmpty(model.Db).
		String("username like ?", "%"+params.Username+"%").
		String("email like ?", "%"+params.Email+"%").
		String("telephone like ?", "%"+params.Telephone+"%").
		String("nickname like ?", "%"+params.Nickname+"%").
		Int64("type=?", params.Type).
		Int64("status=?", params.Status).
		Int64("sex=?", params.Sex).
		RangeTime("birthday between ? and ?", params.Birthday, location).
		RangeTime("updated_at between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)

	// 管理员名称
	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	// 上级用户名称
	if params.ParentName != "" {
		model.Db.AndWhere("parent_id in (" + strings.Join(models.NewUser(nil).FindUserLikeNameIds(params.ParentName), ",") + ")")
	}

	// 国家名称
	if params.CountryName != "" {
		countryIds := models.NewCountry(nil).Field("id").AndWhere("name like ?", "%"+params.CountryName+"%").ColumnString()
		if len(countryIds) == 0 {
			countryIds = append(countryIds, "-1")
		}
		model.Db.AndWhere("country_id in (" + strings.Join(countryIds, ",") + ")")
	}

	data := make([]*indexData, 0)
	model.Field("id,admin_id,parent_id,username,email,nickname,avatar,telephone,sex,birthday,money,freeze_money,type,status,data, INET_NTOA(ip4), updated_at").
		Query(func(rows *sql.Rows) {
			tmp := new(indexData)
			_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.ParentId, &tmp.UserName, &tmp.Email, &tmp.Nickname,
				&tmp.Avatar, &tmp.Telephone, &tmp.Sex, &tmp.Birthday, &tmp.Money, &tmp.FreezeMoney, &tmp.Type, &tmp.Status, &tmp.Data, &tmp.Ip4, &tmp.UpdatedAt)

			// 邀请人用户名
			if tmp.ParentId > 0 {
				parentUserInfo := models.NewUser(nil).AndWhere("id=?", tmp.ParentId).FindOne()
				if parentUserInfo != nil {
					tmp.ParentName = parentUserInfo.UserName
				}
			}

			// 管理员名称
			if tmp.AdminId > 0 {
				adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
				if adminInfo != nil {
					tmp.AdminName = adminInfo.UserName
				}
			}

			//	更新IP4地址
			ip2location, _ := models.GetIp2Location(tmp.Ip4)
			if ip2location != nil {
				tmp.Ip4 = ip2location.Country_long + "." + ip2location.Region + "." + ip2location.City
			}
			//	邀请码
			tmp.InviteCode = models.NewUserInvite(nil).GetInviteCode(tmp.AdminId, tmp.Id)
			data = append(data, tmp)
		})

	body.SuccessJSON(w, body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
