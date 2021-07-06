package main

import (
	"cron_worker/config"
	"cron_worker/cron"
	"time"
)

func main() {
	config.InitConfig()
	cron.InitExecutor()
	cron.InitScheduler()
	cron.InitEtcd()
	for {
		time.Sleep(1 * time.Second)
	}
}
