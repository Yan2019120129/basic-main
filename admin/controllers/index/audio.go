package index

import (
	"basic/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/router"
	"github.com/so68/zfeng/utils/body"
)

type audioData struct {
	Title  string `json:"title"`
	Source string `json:"source"`
}

// Audio 音源
func Audio(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)

	// 订单声音
	orderNums := models.NewProductOrder(nil).AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("status=?", models.ProductOrderStatusPending).Count()
	if orderNums > 0 {
		body.SuccessJSON(w, &audioData{Title: fmt.Sprintf("您有 %v 订单未处理", orderNums), Source: "/assets/mp3/trade.mp3"})
		return
	}

	// 充值声音
	depositNums := models.NewUserWalletOrder(nil).AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("type=?", models.WalletOrderTypeDeposit).AndWhere("status=?", models.WalletOrderStatusPending).Count()
	if depositNums > 0 {
		body.SuccessJSON(w, &audioData{Title: fmt.Sprintf("您有 %v 充值订单未处理", depositNums), Source: "/assets/mp3/deposit.mp3"})
		return
	}

	// 提现声音
	withdrawNums := models.NewUserWalletOrder(nil).AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("type=?", models.WalletOrderTypeWithdraw).AndWhere("status=?", models.WalletOrderStatusPending).Count()
	if withdrawNums > 0 {
		body.SuccessJSON(w, &audioData{Title: fmt.Sprintf("您有 %v 提现订单未处理", withdrawNums), Source: "/assets/mp3/withdraw.mp3"})
		return
	}

	body.SuccessJSON(w, &audioData{Title: "", Source: ""})
}
