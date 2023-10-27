package product

import (
	"basic/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/cache"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/locales"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
)

type createParams struct {
	Id int64 `json:"id" validate:"required"`
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	rds := cache.RedisPool.Get()
	defer rds.Close()

	acceptLanguage := r.Header.Get("Accept-Language")
	claims := router.TokenManager.GetContextClaims(r)
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(claims.AdminId)

	productModel := models.NewProduct(nil)
	productModel.AndWhere("id=?", params.Id).AndWhere("admin_id=?", settingAdminId).AndWhere("status>?", models.ProductStatusDisabled)
	productInfo := productModel.FindOne()
	if productInfo == nil {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "noPermissionToExecute"), -1)
		return
	}

	// 判断用户是否限购
	if productInfo.Nums > 0 {
		userProductOrderNums := models.NewProductOrder(nil).
			AndWhere("product_id=?", productInfo.Id).AndWhere("admin_id=?", claims.AdminId).
			AndWhere("user_id=?", claims.UserId).AndWhere("status=?", models.ProductOrderStatusPending).Count()
		if userProductOrderNums >= productInfo.Nums {
			errMsg := locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "restrictedProducts")
			errMsg = strings.ReplaceAll(errMsg, "{nums}", strconv.FormatInt(productInfo.Nums, 10))
			body.ErrorJSON(w, errMsg, -1)
			return
		}
	}

	// 判断是否达到卖出数量
	if productInfo.Total > 0 && productInfo.Used >= productInfo.Total {
		body.ErrorJSON(w, locales.Manager.GetAdminLocales(rds, settingAdminId, acceptLanguage, "theCurrentProductIsSold"), -1)
		return
	}

	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	nums := 1
	nowTime := time.Now()

	_, err = models.NewProductOrder(tx).
		Field("admin_id", "user_id", "product_id", "order_sn", "money", "nums", "data", "expired_at", "updated_at", "created_at").
		Args(claims.AdminId, claims.UserId, productInfo.Id, utils.NewRandom().OrderSn(), productInfo.Money, nums, productInfo.Data, nowTime.Unix(), nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	//	订单提示音
	//	TODO...
	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
