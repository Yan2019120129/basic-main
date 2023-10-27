package main

import (
	"basic/admin/controllers"
	"basic/models"

	"github.com/so68/zfeng"
)

func main() {
	// 启动后台接口
	adminApp := zfeng.NewApp("./")                                   //	初始化
	adminApp.SetRouteHandle(controllers.Router())                    //	载入后台路由
	adminApp.InitializationAdminRole(models.GetAdminRolesRouter())   //	初始化管理路由
	adminApp.InitializationTokenParams(models.GetAdminTokenParams()) //	初始化管理Token参数
	adminApp.InitializationLocales(models.GetAdminLocales())         // 初始化本地语言 - 缓存在 redis中，所有前台可以不需要初始化
	adminApp.SetCallbackAccessFunc(models.RouterAccessAdminFunc)     //	设置访问日志
	adminApp.ServeFiles("assets")                                    //	资源文件
	adminApp.ListenAndServe("0.0.0.0:8001")                          //	启动监听
}
