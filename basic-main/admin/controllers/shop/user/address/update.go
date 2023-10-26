package address

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
	Id              int64  `json:"id" validate:"required"` //关注ID
	Name            string `json:"name"`                   //收货人名
	Phone           string `json:"phone"`                  //电话
	Country         string `json:"country"`                //国家
	ShippingAddress string `json:"shipping_address"`       //收货地址
	DoorNumber      int64  `json:"door_number"`            //门牌号
	ZipCode         int64  `json:"zip_code"`               //邮编
	IsDefault       int64  `json:"is_default"`             //是否默认：1是，-1否
	Time            int64  `json:"time"`                   //时间
	Status          int64  `json:"status"`                 //状态：1在使用，-1已删除
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
	model := models.NewUserAddress(nil)

	// 开启事务
	tx := database.DbPool.GetTx()
	defer tx.Rollback()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("name=?", params.Name).
		String("phone=?", params.Phone).
		String("country=?", params.Country).
		String("shipping_address=?", params.ShippingAddress).
		Int64("door_number=?", params.DoorNumber).
		Int64("zip_code=?", params.ZipCode).
		Int64("is_default=?", params.IsDefault).
		Int64("time=?", params.Time).
		Int64("status=?", params.Status)
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
