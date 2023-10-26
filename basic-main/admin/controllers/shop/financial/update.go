package financial

import (
	"basic/models"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
	"github.com/so68/zfeng/validator"
	"net/http"
	"strings"
)

type updateParams struct {
	Id               int64  `json:"id"`                 //订单ID
	Quantity         int64  `json:"quantity"`           //数量
	PaymentMethod    string `json:"payment_method"`     //支付方式
	ShippingAddress  string `json:"shipping_address"`   //收货地址
	UserOperate      int64  `json:"user_operate"`       //用户操作：1已下单，2取消订单
	OrderStatus      int64  `json:"order_status"`       //订单状态：0待处理，1已支付，2已发货，3已送达，-1已取消，-2已删除
	UserOperateTime  int64  `json:"user_operate_time"`  //用户操作时间
	StoreOperate     int64  `json:"store_operate"`      //店家操作：1接受订单操作，2拒绝订单操作，3发货操作，4取消订单操作
	StoreOperateTime int64  `json:"store_operate_time"` //店家操作时间
	AdminOperateTime int64  `json:"admin_operate_time"` //管理员操作时间
}

func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updateParams)
	_ = body.ReadJSON(r, params)
	//  参数验证
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 从token获取管理员ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	model := models.NewCommodityOrder(nil)

	// 开启事务
	tx := database.DbPool.GetTx()
	defer tx.Rollback()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		Int64("quantity=?", params.Quantity).
		String("payment_method=?", params.PaymentMethod).
		String("shipping_address=?", params.ShippingAddress).
		Int64("user_operate=?", params.UserOperate).
		Int64("order_status=?", params.OrderStatus).
		Int64("user_operate_time=?", params.UserOperateTime).
		Int64("store_operate=?", params.StoreOperate).
		Int64("store_operate_time=?", params.StoreOperateTime).
		Int64("admin_operate_time=?", params.AdminOperateTime)

	// 是否是超级管理员，不是没有权限修改
	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id = ?", params.Id).Update()
	fmt.Println("err", err)
	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
