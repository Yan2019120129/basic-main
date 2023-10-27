package index

import (
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/utils"
	"github.com/so68/zfeng/utils/body"
)

const (
	UploadFilePath = "/assets/uploads" //	上传文件路径
	MaxUploadSize  = 5 * 1024 * 1000   //	最大上传大小
)

// Upload 文件上传
func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	f, err := os.OpenFile(dir+saveFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	_, _ = io.Copy(f, file)
	body.SuccessJSON(w, saveFilePath)
}
