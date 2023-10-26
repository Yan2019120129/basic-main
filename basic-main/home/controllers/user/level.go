package user

import (
	"basic/models"
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type levelData struct {
	Items   []*levelItem `json:"items"`
	BuyMode string       `json:"buyMode"`
}

type levelItem struct {
	Id    int64   `json:"id"`
	Name  string  `json:"name"`
	Icon  string  `json:"icon"`
	Level int64   `json:"level"`
	Money float64 `json:"money"`
	Days  int64   `json:"days"`
	Data  string  `json:"data"`
}

// Level 用户等级列表
func Level(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := router.TokenManager.GetContextClaims(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)
	settingAdminList := models.NewAdminSetting(nil).GetAdminFieldAllString(settingAdminId)

	acceptLanguage := r.Header.Get("Accept-Language")
	rds := cache.RedisPool.Get()
	defer rds.Close()

	items := make([]*levelItem, 0)
	models.NewUserLevel(nil).
		Field("id", "name", "icon", "level", "money", "days", "data").
		AndWhere("admin_id=?", settingAdminId).AndWhere("status=?", models.UserLevelStatusActivate).
		OrderBy("level asc").Query(func(rows *sql.Rows) {
		temp := new(levelItem)
		_ = rows.Scan(&temp.Id, &temp.Name, &temp.Icon, &temp.Level, &temp.Money, &temp.Days, &temp.Data)
		temp.Name = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, temp.Name)
		temp.Data = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, temp.Data)
		items = append(items, temp)
	})

	body.SuccessJSON(w, &levelData{
		Items:   items,
		BuyMode: settingAdminList["buy_level_mode"],
	})
}
