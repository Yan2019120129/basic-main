package commodity

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/utils/body"
	"net/http"
)

type createParams struct {
	AdminId int64  `json:"product_id" validate:"required,gt=0"`
	Id      int64  `json:"id" validate:"required,gt=0"` //类目ID
	Image   string `json:"image"`                       //类目图片
	Name    string `json:"name"`                        //种类名称
	Date    int64  `json:"date"`                        //时间
	Status  int64  `json:"status"`                      //状态 -2删除 -1禁用 10启用
	Operate int64  `json:"operate"`                     //操作
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)
	tx := database.DbPool.GetTx()
	defer tx.Rollback()
	model := models.NewUserFavorites(nil)
	model.AndWhere("id=?", params.AdminId).FindOne()

	//	TODO... 消费账单...
	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
