package main

import (
	"cron_worker/config"
	"cron_worker/scheduler"
	"time"
)

func main() {
	//t := time.NewTimer(time.Second * 10)
	//for {
	//	<-t.C
	//	fmt.Println("timer running...")
	//	// 需要重置Reset 使 t 重新开始计时
	//	t.Reset(time.Second * 2)
	//}
	//return
	config.InitConfig()
	scheduler.InitScheduler()
	scheduler.InitEtcd()
	//正常退出
	for {
		time.Sleep(1 * time.Second)
	}
}
