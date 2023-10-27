package comment

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/utils/body"
	"net/http"
)

type createParams struct {
	AdminId int64 `json:"product_id" validate:"required"`
	Nums    int64 `json:"nums" validate:"required"`
	Type    int64 `json:"type" validate:"required"`
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
