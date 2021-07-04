package yc_worker

import (
	"cron_worker/config"
	"cron_worker/etcd"
)

func main() {
	config.InitConfig()
	InitScheduler()
	etcd.InitEtcd()
}
