package crontab

import (
	"fmt"
	"time"
)

func InitializeCrontab() {

}

// InitializeTemplate 初始化定时模版
func InitializeTemplate(second time.Duration) {
	ch := time.NewTicker(second)

	for {

		fmt.Println("定时任务模版....")

		<-ch.C
	}
}
