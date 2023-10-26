package role

import (
	"basic/models"
	"database/sql"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
)

type indexParams struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Pagination  *define.Pagination `json:"pagination"`
}

type indexItems struct {
	Name        string          `json:"name"`
	Permissions string          `json:"permissions"`
	Authority   map[string]bool `json:"authority"`
	Description string          `json:"description"`
	UpdatedAt   int64           `json:"updated_at"`
	CreatedAt   int64           `json:"created_at"`
}

// Index 角色列表
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	model := models.NewAdminAuthItem(nil)
	//  模型增加where条件并分页
	define.NewFilterEmpty(model.Db).
		String("name like ?", "%"+params.Name+"%").
		String("description like ?", "%"+params.Description+"%").
		Pagination(params.Pagination)

	data := make([]*indexItems, 0)
	model.AndWhere("type=?", models.AdminAuthItemTypeManage).
		Field("name", "description", "created_at", "updated_at").
		Query(func(rows *sql.Rows) {
			tmp := new(indexItems)
			_ = rows.Scan(&tmp.Name, &tmp.Description, &tmp.CreatedAt, &tmp.UpdatedAt)
			rolesRouteList := utils.GetMapKeys(models.NewAdminAuthChild(nil).GetRolesRouteList([]string{tmp.Name}))
			tmp.Permissions = strings.Join(rolesRouteList, ",")
			tmp.Authority = models.NewAdminAuthChild(nil).GetRouteRoleCheckedList(rolesRouteList)
			data = append(data, tmp)
		})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
