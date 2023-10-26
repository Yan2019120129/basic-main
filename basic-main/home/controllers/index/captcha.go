package index

import (
	"net/http"
	"strconv"

	"github.com/dchest/captcha"
	"github.com/julienschmidt/httprouter"
	"github.com/so68/zfeng/utils/body"
)

// GenerateCaptcha 生成验证码
func GenerateCaptcha(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body.SuccessJSON(w, captcha.NewLen(4))
}

// ImageCaptcha 显示验证码
func ImageCaptcha(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	id := r.URL.Query().Get("id")
	width, _ := strconv.Atoi(r.URL.Query().Get("w"))
	height, _ := strconv.Atoi(r.URL.Query().Get("h"))

	if id == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if width == 0 {
		width = 120
	}

	if height == 0 {
		height = 32
	}
	w.Header().Set("Content-Type", "image/png")
	_ = captcha.WriteImage(w, id, width, height)
}
