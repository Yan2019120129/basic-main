package main

import (
	"basic/home/controllers"
	"basic/models"

	"github.com/so68/zfeng"
)

func main() {
	// 启动前台接口
	homeApp := zfeng.NewApp("./")                                  //	初始化
	homeApp.SetRouteHandle(controllers.Router())                   //	载入前台路由
	homeApp.InitializationTokenParams(models.GetHomeTokenParams()) //	前台Token参数
	homeApp.SetCallbackAccessFunc(models.RouterAccessHomeFunc)     //	设置访问日志
	homeApp.ServeFiles("assets")                                   //	资源文件
	homeApp.ListenAndServe("0.0.0.0:8004")                         //	启动监听
}
