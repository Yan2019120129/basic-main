package index

import (
	"basic/tools/utils"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/utils/body"
)

type tablesData struct {
	Name    string `json:"name"`    //	表名
	Comment string `json:"comment"` //	注释
}

type columnsParams struct {
	Name string `json:"name"` //	表名
}

// Tables 获取所有表名
func Tables(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rows, err := database.DbPool.GetDb().Query("select table_name,table_comment from information_schema.TABLES where table_schema = '" + zfeng.ConfInfo.Database.Dbname + "'")
	if err != nil {
		panic(err)
	}
	var data []*tablesData
	for rows.Next() {
		tmp := new(tablesData)
		err = rows.Scan(&tmp.Name, &tmp.Comment)
		if err != nil {
			panic(err)
		}
		data = append(data, tmp)
	}

	body.SuccessJSON(w, data)
}

// Columns 表字段列
func Columns(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(columnsParams)
	_ = body.ReadJSON(r, params)

	body.SuccessJSON(w, utils.NewOperate().Columns(params.Name))
}
