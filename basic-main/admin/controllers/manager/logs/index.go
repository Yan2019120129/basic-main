package logs

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
	AdminName  string                 `json:"admin_name"`
	UserName   string                 `json:"username"`
	UserAgent  string                 `json:"user_agent"`
	Lang       string                 `json:"lang"`
	Name       string                 `json:"name"`
	Type       int64                  `json:"type"`
	Ip4        string                 `json:"ip4"`
	Route      string                 `json:"route"`
	Data       string                 `json:"data"`
	CreatedAt  *define.RangeTimeParam `json:"created_at"`
	Pagination *define.Pagination     `json:"pagination"` //	分页
}

type indexItems struct {
	AdminName string `json:"admin_name"` //	管理名称
	UserName  string `json:"username"`   //	用户名
	models.AccessLogsAttrs
}

// Index 管理员配置 列表
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//  获取参数
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取请求管理员的uid
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)

	//  实例化模型
	model := models.NewAccessLogs(nil)
	model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	//  添加where条件 和分页
	define.NewFilterEmpty(model.Db).
		String("name like ?", "%"+params.Name+"%").
		String("INET_NTOA(ip4) like ?", "%"+params.Ip4+"%").
		String("route like ?", "%"+params.Route+"%").
		String("data like ?", "%"+params.Data+"%").
		String("user_agent like ?", "%"+params.UserAgent+"%").
		String("lang like ?", "%"+params.Lang+"%").
		Int64("type=?", params.Type).
		RangeTime("created_at between ? and ?", params.CreatedAt, location).
		Pagination(params.Pagination)

	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	if params.UserName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewUser(nil).FindUserLikeNameIds(params.UserName), ",") + ")")
	}

	//  查询
	data := make([]*indexItems, 0)
	model.Field("id", "admin_id", "user_id", "name", "type", "INET_NTOA(ip4)", "route", "user_agent", "lang", "data", "created_at").
		Query(func(rows *sql.Rows) {
			tmp := new(indexItems)
			_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.UserId, &tmp.Name, &tmp.Type, &tmp.Ip4, &tmp.Route, &tmp.UserAgent, &tmp.Lang, &tmp.Data, &tmp.CreatedAt)
			adminInfo := models.NewAdminUser(nil).AndWhere("id=?", tmp.AdminId).FindOne()
			if adminInfo != nil {
				tmp.AdminName = adminInfo.UserName
			}

			userInfo := models.NewUser(nil).AndWhere("id=?", tmp.UserId).FindOne()
			if userInfo != nil {
				tmp.UserName = userInfo.UserName
			}

			//	更新IP4地址
			ip2location, _ := models.GetIp2Location(tmp.Ip4)
			if ip2location != nil {
				tmp.Ip4 = ip2location.Country_long + "." + ip2location.Region + "." + ip2location.City + "[" + tmp.Ip4 + "]"
			}

			data = append(data, tmp)
		})

	//  统计数量
	count := model.Count()
	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: count,
	})
}
