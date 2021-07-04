package api

import (
	"cron_master/utils"
	"github.com/gin-gonic/gin"
)

func GetSysInfo(c *gin.Context) {
	utils.OkWithData(utils.SystemInfo(), c)
}
