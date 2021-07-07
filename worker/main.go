package main

import (
	"cron_worker/config"
	"cron_worker/cron"
	"cron_worker/model"
	"time"
)

func main() {
	config.InitConfig()
	model.InitModels()
	cron.InitExecutor()
	cron.InitScheduler()
	cron.InitEtcd()
	for {
		time.Sleep(1 * time.Second)
	}
}
