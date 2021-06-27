package main

import (
	"yc_cron/config"
	"yc_cron/utils"
)

func main() {
	config.InitConfig()
	utils.StartServer()
}
