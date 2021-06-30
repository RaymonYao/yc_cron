package api

import (
	"github.com/gin-gonic/gin"
	"yc_cron/utils"
)

func GetSysInfo(c *gin.Context) {
	utils.OkWithData(utils.SystemInfo(), c)
}
