package models

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/logs"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"go.uber.org/zap"
)

const (
	AccessLogsTypeAdmin = 1 //	后端日志
	AccessLogsTypeHome  = 2 //	前端日志
)

type AccessLogsAttrs struct {
	Id        int64  `json:"id"`         //主键
	AdminId   int64  `json:"admin_id"`   //管理员ID
	UserId    int64  `json:"user_id"`    //用户ID
	Type      int64  `json:"type"`       //日志类型
	Name      string `json:"name"`       //标题
	Ip4       string `json:"ip4"`        //IP4地址
	UserAgent string `json:"user_agent"` //ua信息
	Lang      string `json:"lang"`       //语言信息
	Route     string `json:"route"`      //操作路由
	Data      string `json:"data"`       //数据
	CreatedAt int64  `json:"created_at"` //时间
}

type AccessLogs struct {
	define.Db
}

func NewAccessLogs(tx *sql.Tx) *AccessLogs {
	return &AccessLogs{
		database.DbPool.NewDb(tx).Table("access_logs"),
	}
}

// AndWhere where条件
func (c *AccessLogs) AndWhere(str string, arg ...any) *AccessLogs {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *AccessLogs) FindOne() *AccessLogsAttrs {
	attrs := new(AccessLogsAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.UserId, &attrs.Type, &attrs.Name, &attrs.Ip4, &attrs.UserAgent, &attrs.Lang, &attrs.Route, &attrs.Data, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *AccessLogs) FindMany() []*AccessLogsAttrs {
	data := make([]*AccessLogsAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(AccessLogsAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.UserId, &tmp.Type, &tmp.Name, &tmp.Ip4, &tmp.UserAgent, &tmp.Lang, &tmp.Route, &tmp.Data, &tmp.CreatedAt)
		data = append(data, tmp)
	})
	return data
}

// UniqueVisitor 网站独立访客
func (c *AccessLogs) UniqueVisitor(adminIds []string, betweenTime []int64) int64 {
	c.Field("count(DISTINCT ip4)").Where("AND", "admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("type=?", AccessLogsTypeHome).AndWhere("user_id>0")
	if len(betweenTime) == 2 {
		c.AndWhere("created_at between ? and ?", betweenTime[0], betweenTime[1])
	}
	return c.Count()
}

// RouterAccessAdminFunc 路由管理日志方法
func RouterAccessAdminFunc(handleParams *router.Handle, r *http.Request, claims *router.Claims) {
	RouteAccessFunc(AccessLogsTypeAdmin, handleParams, r, claims)
}

// RouterAccessHomeFunc 路由前台日志方法
func RouterAccessHomeFunc(handleParams *router.Handle, r *http.Request, claims *router.Claims) {
	RouteAccessFunc(AccessLogsTypeHome, handleParams, r, claims)
}

// RouteAccessFunc 路由日志方法
func RouteAccessFunc(routerId int64, handleParams *router.Handle, r *http.Request, claims *router.Claims) {
	//	验证的路由， 没有验证的路由， 后端跟前端
	var adminId, userId int64
	if claims != nil {
		adminId = claims.AdminId
		userId = claims.UserId
	}

	//	所有没有验证的方法,  没有管理ID,
	if adminId == 0 {
		adminId = NewAdminUser(nil).GetDomainAdminId(r)
	}

	//	过滤操作 客服会话｜提示声音｜websocket链接｜登陆｜注册|更新密码|管理员更新 | 用户修改登录密码｜用户修改安全密码
	if handleParams.Route == "/chat/conversation" ||
		handleParams.Route == "/login" || handleParams.Route == "/register" ||
		handleParams.Route == "/update/password" || handleParams.Route == "/manage/update" ||
		handleParams.Route == "/user/update/password" || handleParams.Route == "/user/update/security" ||
		handleParams.Route == "/audio" || handleParams.Route == "/chat/ws" {
		return
	}

	data := `{"GET": ` + r.URL.Query().Encode() + `, "POST": ` + body.GetBody(r) + `}`
	if strings.Contains(handleParams.Route, "login") {
		data = ""
	}
	nowTime := time.Now().Unix()

	logs.Logger.Debug("access", zap.String("name", handleParams.Name), zap.String("method", handleParams.Method), zap.String("router", handleParams.Route), zap.String("data", data))
	_, _ = NewAccessLogs(nil).
		Field("admin_id", "user_id", "type", "name", "ip4", "user_agent", "lang", "route", "data", "created_at").
		Value("?", "?", "?", "?", "INET_ATON(?)", "?", "?", "?", "?", "?").
		Args(adminId, userId, routerId, handleParams.Name, utils.GetUserRealIP(r), r.Header.Get("User-Agent"), r.Header.Get("Accept-Language"), handleParams.Route, data, nowTime).
		Insert()
}
