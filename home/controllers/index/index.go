package index

import (
	"basic/models"
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type categoryItem struct {
	Id    int64  `json:"id"`
	Image string `json:"image"`
	Name  string `json:"name"`
}

type productItem struct {
	Id     int64   `json:"id"`
	Images []any   `json:"images"`
	Name   string  `json:"name"`
	Money  float64 `json:"money"`
	Type   int64   `json:"type"`
	Sales  int64   `json:"sales"`
	Nums   int64   `json:"nums"`
	Used   int64   `json:"used"`
	Total  int64   `json:"total"`
	Data   string  `json:"data"`
	Desc   string  `json:"desc"`
}

type AwardItem struct {
	Avatar   string  `json:"avatar"`   //	头像
	UserName string  `json:"username"` //	用户名
	Name     string  `json:"name"`     //	名称
	Money    float64 `json:"money"`    //	金额
}

type indexData struct {
	Banner       []map[string]any `json:"banner"`       //	首页Banner图
	Introduce    string           `json:"introduce"`    //	首页介绍
	Notice       string           `json:"notice"`       //	首页公告
	AwardList    []*AwardItem     `json:"awardList"`    //	奖励公告
	CategoryList []*categoryItem  `json:"categoryList"` //	推荐分类列表
	ProductList  []*productItem   `json:"productList"`  //	推荐产品列表
}

// Index 首页信息
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := models.NewAdminUser(nil).GetDomainAdminId(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	acceptLanguage := r.Header.Get("Accept-Language")
	rds := cache.RedisPool.Get()
	defer rds.Close()

	//	推荐产品，和推荐分类
	categoryList := make([]*categoryItem, 0)
	productList := make([]*productItem, 0)

	//	推荐分类
	models.NewProductCategory(nil).Field("id", "image", "name").
		AndWhere("admin_id=?", settingAdminId).AndWhere("status>?", models.ProductCategoryStatusDisabled).AndWhere("recommend=?", models.ProductCategoryRecommend).
		OrderBy("sort asc").Query(func(rows *sql.Rows) {
		tmp := new(categoryItem)
		_ = rows.Scan(&tmp.Id, &tmp.Image, &tmp.Name)
		tmp.Name = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, tmp.Name)
		categoryList = append(categoryList, tmp)
	})

	//	推荐产品
	models.NewProduct(nil).Field("id", "images", "name", "money", "type", "sales", "nums", "used", "total", "data", "describes").
		AndWhere("admin_id=?", settingAdminId).AndWhere("status>?", models.ProductStatusDisabled).AndWhere("recommend=?", models.ProductRecommend).
		OrderBy("sort asc").Query(func(rows *sql.Rows) {
		tmp := new(productItem)
		var productImages string

		_ = rows.Scan(&tmp.Id, &productImages, &tmp.Name, &tmp.Money, &tmp.Type, &tmp.Sales, &tmp.Nums, &tmp.Used, &tmp.Total, &tmp.Data, &tmp.Desc)
		_ = json.Unmarshal([]byte(productImages), &tmp.Images)

		tmp.Name = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, tmp.Name)
		tmp.Desc = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, tmp.Desc)
		productList = append(productList, tmp)
	})

	//	获取配置参数
	adminSettingList := models.NewAdminSetting(nil).GetAdminFieldAllString(settingAdminId)

	//	如果用户是登陆状态，并且设置了公告跟介绍，那么显示对应的数据
	claims := router.TokenManager.GetHeaderClaims(rds, r)
	if claims != nil {
		userIntroduce := models.NewUserSetting(nil).GetUserSettingValue(claims.UserId, "introduce")
		if userIntroduce != "" {
			adminSettingList["home_introduce"] = userIntroduce
		}

		userNotice := models.NewUserSetting(nil).GetUserSettingValue(claims.UserId, "notice")
		if userNotice != "" {
			adminSettingList["home_notice"] = userNotice
		}
	}

	//	奖励公告
	awardList := make([]*AwardItem, 0)
	_ = json.Unmarshal([]byte(adminSettingList["home_scroll"]), &awardList)

	for _, awardItem := range awardList {
		awardItem.Name = locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, awardItem.Name)
	}
	//	随机
	rand.Shuffle(len(awardList), func(i, j int) { awardList[i], awardList[j] = awardList[j], awardList[i] })

	bannerList := []map[string]any{}
	_ = json.Unmarshal([]byte(adminSettingList["home_banner"]), &bannerList)

	body.SuccessJSON(w, &indexData{
		Banner:       bannerList,
		Introduce:    locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, adminSettingList["home_introduce"]),
		Notice:       locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, adminSettingList["home_notice"]),
		AwardList:    awardList,
		CategoryList: categoryList,
		ProductList:  productList,
	})
}
