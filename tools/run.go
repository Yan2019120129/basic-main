package main

import (
	"basic/admin/controllers"
	"basic/tools/sql"
	"basic/tools/utils"
	"fmt"
	"os"

	"github.com/so68/zfeng"
)

func main() {
	zfeng.NewApp("./")

	for {
		fmt.Println("1) 重置数据库")
		fmt.Println("2) 生成模型")
		fmt.Println("0) 退出操作")

		var opt int
		_, _ = fmt.Scan(&opt)

		switch opt {
		case 0:
			os.Exit(0)
		case 1:
			utils.NewOperate().Reload(sql.InitializeTables()).LoadRouterList(controllers.Router(), sql.InitializeAuth())
			fmt.Println("重置数据库成功...")
		case 2:
			fmt.Print("请输入数据库表名:")
			var tableName string
			_, _ = fmt.Scan(&tableName)
			utils.NewOperate().GenerateModel("/models", tableName)
		default:
			fmt.Println("请输入正确的操作标号...")
		}
	}
}
