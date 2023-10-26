package user

import (
	"basic/models"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type teamItem struct {
	Name     string      `json:"name"`
	Items    []*userItem `json:"items"`
	SumMoney float64     `json:"sum_money"`
}

type userItem struct {
	Id        int64   `json:"id"`
	ParentId  int64   `json:"parent_id"`
	Avatar    string  `json:"avatar"`
	UserName  string  `json:"username"`
	Money     float64 `json:"money"`
	CreatedAt int64   `json:"created_at"`
}

// Team 团队列表
func Team(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := router.TokenManager.GetContextClaims(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)
	adminSettingList := models.NewAdminSetting(nil).GetAdminFieldAllString(settingAdminId)
	acceptLanguage := r.Header.Get("Accept-Language")
	rds := cache.RedisPool.Get()
	defer rds.Close()

	//	金字塔等级
	pyramidLevel := models.AdminSettingValueToMapInterfaces(adminSettingList["pyramid_level"])
	var userParentIds []string
	data := make([]*teamItem, 0)
	userParentIds = []string{strconv.FormatInt(claims.UserId, 10)}

	for i := 0; i < len(pyramidLevel); i++ {
		itemTmp := pyramidLevel[i]
		teamName := locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, itemTmp["label"].(string))
		tempTeam := &teamItem{Name: teamName, SumMoney: 0, Items: make([]*userItem, 0)}

		if len(userParentIds) > 0 {
			var userParentIdsTmp []string

			//	用户信息
			models.NewUser(nil).Field("id", "parent_id", "avatar", "username", "created_at").
				AndWhere("parent_id in ("+strings.Join(userParentIds, ",")+")").
				AndWhere("status=?", models.UserStatusActivate).Query(func(rows *sql.Rows) {
				userItemTmp := new(userItem)
				_ = rows.Scan(&userItemTmp.Id, &userItemTmp.ParentId, &userItemTmp.Avatar, &userItemTmp.UserName, &userItemTmp.CreatedAt)

				if userItemTmp.Id > 0 {
					userParentIdsTmp = append(userParentIdsTmp, strconv.FormatInt(userItemTmp.Id, 10))
				}

				//	购买产品金额
				models.NewUserBill(nil).Field("sum(money)").
					AndWhere("user_id=?", userItemTmp.Id).AndWhere("type=?", models.UserBillTypeBuyProduct).QueryRow(func(row *sql.Row) {
					_ = row.Scan(&userItemTmp.Money)
				})
				tempTeam.SumMoney += userItemTmp.Money
				tempTeam.Items = append(tempTeam.Items, userItemTmp)
			})
			userParentIds = userParentIdsTmp
		}
		data = append(data, tempTeam)
	}

	body.SuccessJSON(w, data)
}
