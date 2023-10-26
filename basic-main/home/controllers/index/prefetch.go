package index

import (
	"basic/models"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/utils/body"
)

type country struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Code string `json:"code"`
}

type Language struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
	Icon  string `json:"icon"`
}

type Locale struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type preFetchData struct {
	Config    map[string]interface{} `json:"config"`    //	配置
	Template  map[string]interface{} `json:"template"`  //	模版配置
	Lang      string                 `json:"lang"`      //	当前语言
	LangList  []*Language            `json:"langList"`  //	语言列表
	Locales   []*Locale              `json:"locales"`   //	默认语言包
	Countries []*country             `json:"countries"` //	国家列表
}

// PreFetch 预处理数据
func PreFetch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := models.NewAdminUser(nil).GetDomainAdminId(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	settingAdminList := models.NewAdminSetting(nil).GetAdminFieldAllString(settingAdminId)
	rds := cache.RedisPool.Get()
	defer rds.Close()

	//	获取国家列表
	countries := make([]*country, 0)
	models.NewCountry(nil).Field("id", "alias", "icon", "code").
		AndWhere("admin_id=?", settingAdminId).AndWhere("status>?", models.CountryStatusDisabled).
		OrderBy("sort asc").Query(func(rows *sql.Rows) {
		countryTmp := new(country)
		_ = rows.Scan(&countryTmp.Id, &countryTmp.Name, &countryTmp.Icon, &countryTmp.Code)
		countries = append(countries, countryTmp)
	})

	//	获取管理语言列表
	languageList := make([]*Language, 0)
	models.NewLang(nil).Field("id", "name", "alias", "icon").
		AndWhere("admin_id=?", settingAdminId).AndWhere("status>?", models.LangStatusDisabled).
		OrderBy("sort asc").Query(func(rows *sql.Rows) {
		languageTmp := new(Language)
		_ = rows.Scan(&languageTmp.Id, &languageTmp.Name, &languageTmp.Alias, &languageTmp.Icon)
		languageList = append(languageList, languageTmp)
	})

	//  获取当前语言包
	alias := languageList[0].Alias
	selectLang := models.NewLang(nil).FindAlias(settingAdminId, r.Header.Get("Accept-Language"))
	if selectLang != nil {
		alias = selectLang.Alias
	}
	localesList := make([]*Locale, 0)
	models.NewLangDictionary(nil).Field("field", "value").
		AndWhere("alias=?", alias).
		AndWhere("admin_id=?", settingAdminId).AndWhere("type=?", models.LangDictionaryTypeHomeTranslate).
		Query(func(rows *sql.Rows) {
			localeTmp := new(Locale)
			_ = rows.Scan(&localeTmp.Label, &localeTmp.Value)
			localesList = append(localesList, localeTmp)
		})

	//	获取主题模版
	adminInfo := models.NewAdminUser(nil).AndWhere("id=?", settingAdminId).FindOne()
	adminExtra := new(models.AdminUserExtraAttrs)
	_ = json.Unmarshal([]byte(adminInfo.Extra), &adminExtra)
	if adminExtra.Template == "" {
		adminExtra.Template = "default"
	}

	body.SuccessJSON(w, &preFetchData{
		Config: map[string]interface{}{
			"site_logo":   settingAdminList["site_logo"],                                           //	站点Logo
			"site_name":   settingAdminList["site_name"],                                           //	站点名称
			"home_online": settingAdminList["home_online"],                                         //	客服链接
			"onlineIcon":  settingAdminList["onlineIcon"],                                          //	客服图标
			"admin_tabs":  models.AdminSettingValueToMapInterfaces(settingAdminList["admin_tabs"]), //	Tabs导航
		},
		Template: map[string]interface{}{
			"template":          adminExtra.Template,                                                           //	主题模版
			"color_primary":     settingAdminList["color_primary"],                                             //	主题色
			"color_secondary":   settingAdminList["color_secondary"],                                           //	辅助色
			"color_accent":      settingAdminList["color_accent"],                                              //	强调色
			"template_basic":    models.AdminSettingValueToMapInterface(settingAdminList["template_basic"]),    //	模版基础配置
			"template_show":     models.AdminSettingValueToMapInterface(settingAdminList["template_show"]),     //	模版显示配置
			"template_login":    models.AdminSettingValueToMapInterface(settingAdminList["template_login"]),    //	模版登陆配置
			"template_register": models.AdminSettingValueToMapInterface(settingAdminList["template_register"]), //	模版注册配置
			"template_wallet":   models.AdminSettingValueToMapInterface(settingAdminList["template_wallet"]),   //	模版钱包配置
			"template_verify":   models.AdminSettingValueToMapInterface(settingAdminList["template_verify"]),   //	模版验证配置
		},
		Lang:      alias,
		LangList:  languageList,
		Locales:   localesList,
		Countries: countries,
	})
}
