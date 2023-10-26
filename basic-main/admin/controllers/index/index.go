package index

import (
	"basic/models"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type statistics struct {
	Name      string `json:"name"`      //	名称
	Icon      string `json:"icon"`      //	图标
	Color     string `json:"color"`     //	背景颜色
	Today     any    `json:"today"`     //	今日
	Yesterday any    `json:"yesterday"` //	昨日
	Total     any    `json:"total"`     //	总数
}

type series struct {
	Name   string `json:"name"`   //	名称
	Type   string `json:"type"`   //	线类型
	Smooth bool   `json:"smooth"` // 	平滑
	Data   []any  `json:"data"`   //	数据
}

type echarts struct {
	Category   []string  `json:"category"` //	日期
	SeriesList []*series `json:"series"`   //	数据
}

type indexData struct {
	Items       [][]*statistics `json:"items"`      //	图标信息
	EchartsList *echarts        `json:"echarts"`    //	图标数据
	OnlineNums  int64           `json:"onlineNums"` //	在线人数
	UnreadNums  int64           `json:"unreadNums"` //	未读消息
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)

	//	获取管理员时区
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))
	nowTime := time.Now()
	nowTime1 := time.Now().AddDate(0, 0, -1)
	nowTime0Unix := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, location)
	nowTime1Unix := time.Date(nowTime1.Year(), nowTime1.Month(), nowTime1.Day(), 0, 0, 0, 0, location)

	// 活跃人数	今日人数｜昨日人数｜总人数
	visitorToday := models.NewAccessLogs(nil).UniqueVisitor(adminIds, []int64{nowTime0Unix.Unix(), nowTime.Unix()})
	visitorYesterday := models.NewAccessLogs(nil).UniqueVisitor(adminIds, []int64{nowTime1Unix.Unix(), nowTime0Unix.Unix()})
	visitorSum := models.NewAccessLogs(nil).UniqueVisitor(adminIds, []int64{})

	// 用户数	今日注册数｜昨日注册数｜总数
	userToday := models.NewUser(nil).Field("count(*)").AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("created_at between ? and ?", nowTime0Unix.Unix(), nowTime.Unix()).Count()
	userYesterday := models.NewUser(nil).Field("count(*)").AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("created_at between ? and ?", nowTime1Unix.Unix(), nowTime0Unix.Unix()).Count()
	userSum := models.NewUser(nil).Field("count(*)").AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")").Count()

	// 充值数 	今日充值数｜昨日充值数｜总充值
	var depositToday float64
	var depositYesterday float64
	var depositSum float64
	models.NewUserWalletOrder(nil).Field("sum(money)").AndWhere("user_type=?", models.UserTypeReality).AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("type<?", models.WalletOrderTypeWithdraw).AndWhere("status=?", models.WalletOrderStatusComplete).AndWhere("created_at between ? and ?", nowTime0Unix.Unix(), nowTime.Unix()).QueryRow(func(row *sql.Row) {
		_ = row.Scan(&depositToday)
	})
	models.NewUserWalletOrder(nil).Field("sum(money)").AndWhere("user_type=?", models.UserTypeReality).AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("type<?", models.WalletOrderTypeWithdraw).AndWhere("status=?", models.WalletOrderStatusComplete).AndWhere("created_at between ? and ?", nowTime1Unix.Unix(), nowTime0Unix.Unix()).QueryRow(func(row *sql.Row) {
		_ = row.Scan(&depositYesterday)
	})
	models.NewUserWalletOrder(nil).Field("sum(money)").AndWhere("user_type=?", models.UserTypeReality).AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("type<?", models.WalletOrderTypeWithdraw).AndWhere("status=?", models.WalletOrderStatusComplete).QueryRow(func(row *sql.Row) {
		_ = row.Scan(&depositSum)
	})

	// 提现数	今日提现数｜昨日提现数｜总提现
	var withdrawToday float64
	var withdrawYesterday float64
	var withdrawSum float64
	models.NewUserWalletOrder(nil).Field("sum(money)").AndWhere("user_type=?", models.UserTypeReality).AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("type>=?", models.WalletOrderTypeWithdraw).AndWhere("status=?", models.WalletOrderStatusComplete).AndWhere("created_at between ? and ?", nowTime0Unix.Unix(), nowTime.Unix()).QueryRow(func(row *sql.Row) {
		_ = row.Scan(&withdrawToday)
	})
	models.NewUserWalletOrder(nil).Field("sum(money)").AndWhere("user_type=?", models.UserTypeReality).AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("type>=?", models.WalletOrderTypeWithdraw).AndWhere("status=?", models.WalletOrderStatusComplete).AndWhere("created_at between ? and ?", nowTime1Unix.Unix(), nowTime0Unix.Unix()).QueryRow(func(row *sql.Row) {
		_ = row.Scan(&withdrawYesterday)
	})
	models.NewUserWalletOrder(nil).Field("sum(money)").AndWhere("user_type=?", models.UserTypeReality).AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("type>=?", models.WalletOrderTypeWithdraw).AndWhere("status=?", models.WalletOrderStatusComplete).QueryRow(func(row *sql.Row) {
		_ = row.Scan(&withdrawSum)
	})

	var category []string
	var visitorNumList []any
	var userNumList []any
	var depositNumList []any
	var withdrawNumList []any

	for i := -14; i <= 0; i++ {
		nowTimeTmp := time.Now().AddDate(0, 0, i)
		sourceTime := time.Date(nowTimeTmp.Year(), nowTimeTmp.Month(), nowTimeTmp.Day(), 0, 0, 0, 0, location)
		staTime := sourceTime.Unix()
		category = append(category, sourceTime.Format("01/02"))

		// 	访客数
		visitorNum := models.NewAccessLogs(nil).UniqueVisitor(adminIds, []int64{staTime, staTime + 86399})

		// 	用户数
		userNum := models.NewUser(nil).Field("count(*)").AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("created_at between ? and ?", staTime, staTime+86399).Count()

		//	充值量
		var depositNum float64
		var withdrawNum float64
		models.NewUserWalletOrder(nil).Field("sum(money)").AndWhere("user_type=?", models.UserTypeReality).AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("type<?", models.WalletOrderTypeWithdraw).AndWhere("status=?", models.WalletOrderStatusComplete).AndWhere("created_at between ? and ?", staTime, staTime+86399).QueryRow(func(row *sql.Row) {
			_ = row.Scan(&depositNum)
		})

		// 提现量
		models.NewUserWalletOrder(nil).Field("sum(money)").AndWhere("user_type=?", models.UserTypeReality).AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("type>=?", models.WalletOrderTypeWithdraw).AndWhere("status=?", models.WalletOrderStatusComplete).AndWhere("created_at between ? and ?", staTime, staTime+86399).QueryRow(func(row *sql.Row) {
			_ = row.Scan(&withdrawNum)
		})

		visitorNumList = append(visitorNumList, visitorNum)
		userNumList = append(userNumList, userNum)
		depositNumList = append(depositNumList, depositNum)
		withdrawNumList = append(withdrawNumList, withdrawNum)
	}

	var unreadNums, onlineNums int64
	onlineUserInfo := models.NewUser(nil).AndWhere("admin_id=?", adminId).AndWhere("type=?", models.UserTypeOnline).AndWhere("status=?", models.UserStatusActivate).FindOne()
	if onlineUserInfo != nil {
		unreadNums = models.NewChatConversationMessage(nil).AndWhere("receiver_id=?", onlineUserInfo.Id).AndWhere("unread=?", models.ChatConversationMessageUnread).Count()
	}

	body.SuccessJSON(w, &indexData{
		OnlineNums: onlineNums,
		UnreadNums: unreadNums,
		Items: [][]*statistics{
			{
				&statistics{
					Name:      "访客数",
					Icon:      "sym_o_person",
					Color:     "bg-primary",
					Today:     visitorToday,
					Yesterday: visitorYesterday,
					Total:     visitorSum,
				},
				&statistics{
					Name:      "用户数",
					Icon:      "sym_o_person_add",
					Color:     "bg-secondary",
					Today:     userToday,
					Yesterday: userYesterday,
					Total:     userSum,
				},
				&statistics{
					Name:      "充值量",
					Icon:      "sym_o_credit_card",
					Color:     "bg-accent",
					Today:     depositToday,
					Yesterday: depositYesterday,
					Total:     depositSum,
				},
				&statistics{
					Name:      "提现量",
					Icon:      "sym_o_payments",
					Color:     "bg-dark",
					Today:     withdrawToday,
					Yesterday: withdrawYesterday,
					Total:     withdrawSum,
				},
			},
		},
		EchartsList: &echarts{
			Category: category,
			SeriesList: []*series{
				{
					Name:   "访客数",
					Type:   "line",
					Smooth: true,
					Data:   visitorNumList,
				}, {
					Name:   "用户数",
					Type:   "line",
					Smooth: true,
					Data:   userNumList,
				}, {
					Name:   "充值量",
					Type:   "line",
					Smooth: true,
					Data:   depositNumList,
				}, {
					Name:   "提现量",
					Type:   "line",
					Smooth: true,
					Data:   withdrawNumList,
				},
			},
		},
	})
}
