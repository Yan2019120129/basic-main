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

func Info(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := router.TokenManager.GetContextClaims(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)
	acceptLanguage := r.Header.Get("Accept-Language")
	rds := cache.RedisPool.Get()
	defer rds.Close()

	userInfo := models.NewUser(nil).AndWhere("id=?", claims.UserId).FindOne()
	if userInfo == nil {
		panic("home/controller/user/info.go | 当前用户不存在")
	}

	//	验证信息
	verifyInfo := new(models.UserVerifyInfo)
	models.NewUserVerify(nil).Field("status", "data").AndWhere("user_id=?", userInfo.Id).QueryRow(func(row *sql.Row) {
		_ = row.Scan(&verifyInfo.Status, &verifyInfo.Data)
	})

	//	等级信息
	levelInfo := new(models.UserLevelInfo)
	models.NewUserLevelOrder(nil).Field("level_id").AndWhere("user_id=?", userInfo.Id).AndWhere("status=?", models.UserLevelOrderStatusActivate).QueryRow(func(row *sql.Row) {
		var userLevelId int64
		_ = row.Scan(&userLevelId)

		if userLevelId > 0 {
			models.NewUserLevel(nil).Field("id", "level", "name", "icon", "days", "money", "created_at", "updatedAt").AndWhere("id=?", userLevelId).QueryRow(func(row *sql.Row) {
				_ = row.Scan(&levelInfo.Id, &levelInfo.Level, &levelInfo.Name, &levelInfo.Icon, &levelInfo.Days, &levelInfo.Money, &levelInfo.CreatedAt, &levelInfo.UpdatedAt)
			})
		}
		if levelInfo.Name != "" {
			levelInfo.Name = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, levelInfo.Name)
		}
	})

	body.SuccessJSON(w, &models.UserInfo{
		Id:          userInfo.Id,
		CountryId:   userInfo.CountryId,
		UserName:    userInfo.UserName,
		Nickname:    userInfo.Nickname,
		Email:       userInfo.Email,
		Telephone:   userInfo.Telephone,
		Avatar:      userInfo.Avatar,
		Sex:         userInfo.Sex,
		Birthday:    userInfo.Birthday,
		Money:       userInfo.Money,
		FreezeMoney: userInfo.FreezeMoney,
		Data:        userInfo.Data,
		VerifyInfo:  verifyInfo,
		LevelInfo:   levelInfo,
		InviteCode:  models.NewUserInvite(nil).GetInviteCode(userInfo.AdminId, userInfo.Id),
		CreatedAt:   userInfo.CreatedAt,
		UpdatedAt:   userInfo.UpdatedAt,
	})
}
