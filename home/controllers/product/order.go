package product

import (
	"basic/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type orderData struct {
	Id        int64                  `json:"id"`
	ProductId int64                  `json:"product_id"`
	OrderSn   string                 `json:"order_sn"`
	Money     float64                `json:"money"`
	Nums      int64                  `json:"nums"`
	Data      map[string]interface{} `json:"data"`
	Status    int64                  `json:"status"`
	NowTime   int64                  `json:"now_time"`
	ExpiredAt int64                  `json:"expired_at"` //过期时间
	UpdatedAt int64                  `json:"updated_at"` //更新时间
	CreatedAt int64                  `json:"created_at"` //创建时间
	Product   *detailsData           `json:"product"`
}

func Order(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := router.TokenManager.GetContextClaims(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)

	acceptLanguage := r.Header.Get("Accept-Language")
	rds := cache.RedisPool.Get()
	defer rds.Close()

	data := make([]*orderData, 0)
	nowTime := time.Now()
	models.NewProductOrder(nil).Field("id", "product_id", "order_sn", "money", "nums", "data", "status", "expired_at", "updated_at", "created_at").
		AndWhere("admin_id=?", claims.AdminId).AndWhere("user_id=?", claims.UserId).AndWhere("status=?", models.ProductOrderStatusPending).Query(func(rows *sql.Rows) {
		temp := new(orderData)
		var orderDataStr string
		_ = rows.Scan(&temp.Id, &temp.ProductId, &temp.OrderSn, &temp.Money, &temp.Nums, &orderDataStr, &temp.Status, &temp.ExpiredAt, &temp.UpdatedAt, &temp.CreatedAt)
		if orderDataStr != "" {
			_ = json.Unmarshal([]byte(orderDataStr), &temp.Data)
		}
		temp.NowTime = nowTime.Unix()

		// 产品信息
		productModel := models.NewProduct(nil)
		productModel.AndWhere("id=?", temp.ProductId).AndWhere("admin_id=?", settingAdminId).AndWhere("status>?", models.ProductStatusDelete)
		productInfo := productModel.FindOne()
		if productInfo != nil {
			images := make([]map[string]string, 0)
			if productInfo.Images != "" {
				_ = json.Unmarshal([]byte(productInfo.Images), &images)
			}

			var productData map[string]interface{}
			if productInfo.Data != "" {
				_ = json.Unmarshal([]byte(productInfo.Data), &productData)
			}
			temp.Product = &detailsData{
				Id:        productInfo.Id,
				Name:      locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, productInfo.Name),
				Images:    images,
				Describes: locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, productInfo.Describes),
				Data:      productData,
				Money:     productInfo.Money,
				Recommend: productInfo.Recommend,
				Sales:     productInfo.Sales,
			}

			data = append(data, temp)
		}
	})

	body.SuccessJSON(w, data)
}
