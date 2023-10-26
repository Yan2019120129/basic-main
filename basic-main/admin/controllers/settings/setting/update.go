package setting

import (
	"basic/models"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type settingItem struct {
	Id    int64  `json:"id"`
	Type  string `json:"type"`
	Field string `json:"field"`
	Value any    `json:"value"`
}

// Update 更新设置
func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var params map[string]settingItem
	_ = body.ReadJSON(r, &params)

	rds := cache.RedisPool.Get()
	defer rds.Close()

	adminId := router.TokenManager.GetContextClaims(r).AdminId
	for _, item := range params {
		// 不能修改提示字典语言中修改
		if item.Field == models.AdminSettingSiteName || item.Field == models.AdminSettingIntroduce ||
			item.Field == models.AdminSettingNotice || item.Field == models.AdminSettingPrivacyPolicy ||
			item.Field == models.AdminSettingServiceAgreement || item.Field == models.AdminSettingDepositTip ||
			item.Field == models.AdminSettingWithdrawTip {
			continue
		}

		model := models.NewAdminSetting(nil)

		// 保存缓存Token参数
		if item.Type == models.SettingTypeJson && item.Field == models.UpdateAdminTokenParamsField {
			tokenParams := new(router.TokenParams)
			tokenParamsBytes, _ := json.Marshal(item.Value)
			_ = json.Unmarshal(tokenParamsBytes, &tokenParams)
			router.TokenManager.SetTokenParams(rds, models.TokenParamsPrefix(models.HomePrefixTokenKey, adminId), tokenParams)
		}

		// 处理对象数据
		switch item.Type {
		case models.SettingTypeCheckbox, models.SettingTypeJson:
			valueBytes, _ := json.Marshal(item.Value.(map[string]any))
			item.Value = string(valueBytes)
		case models.SettingTypeImages, models.SettingTypeChildren:
			valueBytes, _ := json.Marshal(item.Value.([]any))
			item.Value = string(valueBytes)
		}

		model.Value("value=?").Args(item.Value)
		if adminId != models.AdminUserSupermanId {
			adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
			model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
		}
		_, err := model.AndWhere("id=?", item.Id).Update()
		if err != nil {
			panic(err)
		}
	}
	body.SuccessJSON(w, "ok")
}
