package menu

import (
	"basic/models"
	"database/sql"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/database/define"
	"github.com/so68/zfeng/utils/body"
)

type indexParams struct {
	Name       string             `json:"name"`
	ParentName string             `json:"parent_name"`
	Status     int64              `json:"status"`
	Pagination *define.Pagination `json:"pagination"` //	分页
}

type indexData struct {
	ParentName string `json:"parent_name"`
	models.AdminMenuAttrs
}

// Index 菜单列表
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	model := models.NewAdminMenu(nil)
	define.NewFilterEmpty(model.Db).
		String("name like ?", "%"+params.Name+"%").
		Int64("status = ?", params.Status).
		Pagination(params.Pagination)

	if params.ParentName != "" {
		parentNameIds := models.NewAdminMenu(nil).Field("id").AndWhere("name like ?", "%"+params.ParentName+"%").ColumnString()
		if len(parentNameIds) == 0 {
			parentNameIds = append(parentNameIds, "-1")
		}
		model.Db.AndWhere("parent in (" + strings.Join(parentNameIds, ",") + ")")
	}

	data := make([]*indexData, 0)
	model.Field("id", "name", "parent", "route", "sort", "status", "data").Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.Name, &tmp.Parent, &tmp.Route, &tmp.Sort, &tmp.Status, &tmp.Data)

		parentInfo := models.NewAdminMenu(nil).AndWhere("id=?", tmp.Parent).FindOne()
		if parentInfo != nil {
			tmp.ParentName = parentInfo.Name
		}
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
