package dictionary

import (
	"basic/models"
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"path"
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
)

const (
	UploadFilePath = "/assets/uploads" //	上传文件路径
	MaxUploadSize  = 5 * 1024 * 1000   //	最大上传大小
)

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	_ = r.ParseMultipartForm(MaxUploadSize)
	file, handler, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//	生成上传文件路径
	dir, _ := os.Getwd()
	filenameWithSuffix := path.Base(handler.Filename)
	saveFileDir := UploadFilePath + "/" + time.Now().Format("200601")
	saveFileName := utils.NewRandom().String(12)
	saveFilePath := saveFileDir + "/" + saveFileName + path.Ext(filenameWithSuffix)
	if !utils.PathExists(dir + saveFileDir) {
		utils.PathMkdirAll(dir + saveFileDir)
	}
	//	保持文件
	f, err := os.OpenFile(dir+saveFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	_, _ = io.Copy(f, file)

	//	读取解析csv文件
	csvFile, err := os.Open(dir + saveFilePath)
	if err != nil {
		panic(err)
	}

	//	读取文件并且, 载入数据库
	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	csvRead := csv.NewReader(csvFile)

	// 管理员所有语言别名
	adminLangAliasList := models.NewLang(nil).Field("alias").AndWhere("admin_id=?", adminId).AndWhere("status>?", models.LangStatusDisabled).ColumnString()
	var insertIntoList []string

	aliasLocales := map[string]map[string]string{}
	i := 0
	for {
		record, err := csvRead.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			body.ErrorJSON(w, err.Error(), -1)
			return
		}

		//	第二行开始插入
		if i > 0 {
			recordDownload := &downloadData{
				Alias: record[0],
				Type:  record[1],
				Name:  record[2],
				Field: record[3],
				Value: record[4],
			}

			// 	如果当前别名已经存在，那么不需要添加
			dictionaryInfo := models.NewLangDictionary(nil).AndWhere("alias=?", recordDownload.Alias).AndWhere("field=?", recordDownload.Field).AndWhere("admin_id=?", adminId).AndWhere("type=?", recordDownload.Type).FindOne()
			if dictionaryInfo != nil {
				continue
			}

			//	判断语言别名是否存在
			if utils.SliceStringIndexOf(recordDownload.Alias, adminLangAliasList) == -1 {
				body.ErrorJSON(w, "不存在的语言别名 【"+recordDownload.Alias+"】", -1)
				return
			}

			if _, ok := aliasLocales[recordDownload.Alias]; !ok {
				aliasLocales[recordDownload.Alias] = map[string]string{}
			}

			aliasLocales[recordDownload.Alias][recordDownload.Field] = recordDownload.Value
			// 	添加插入语句
			insertIntoList = append(insertIntoList, "("+strconv.FormatInt(adminId, 10)+","+recordDownload.Type+",'"+recordDownload.Alias+"', '"+recordDownload.Name+"','"+recordDownload.Field+"', '"+recordDownload.Value+"')")
		}
		i++
	}

	//	新增数据到管理员语言字典
	_, err = tx.Exec("INSERT INTO lang_dictionary(admin_id, type, alias, name, field, value) VALUES " + strings.Join(insertIntoList, ","))
	if err != nil {
		panic(err)
	}

	_ = tx.Commit()

	//	更新系统内存语言对
	rds := cache.RedisPool.Get()
	defer rds.Close()
	for alias, locale := range aliasLocales {
		locales.Manager.SetAdminLocalesAll(rds, adminId, alias, locale)
	}

	body.SuccessJSON(w, "ok")
}
