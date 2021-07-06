package main

import (
	"cron_worker/config"
	"cron_worker/scheduler"
	"time"
)

func main() {
	config.InitConfig()
	scheduler.InitScheduler()
	scheduler.InitEtcd()
	//正常退出
	for {
		time.Sleep(1 * time.Second)
	}
}
