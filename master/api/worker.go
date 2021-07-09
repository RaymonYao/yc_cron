package api

import (
	"cron_master/model"
	"cron_master/model/request"
	"cron_master/model/response"
	"cron_master/service"
	"cron_master/utils"
	"github.com/gin-gonic/gin"
)

func GetWorkerList(c *gin.Context) {
	var (
		pageInfo request.ComPageInfo
		list     []model.Worker
		total    int
		err      error
	)
	_ = c.ShouldBindJSON(&pageInfo)
	if list, total, err = service.GetWorkerList(&pageInfo); err != nil {
		utils.FailWithMessage("获取失败, Message: "+err.Error(), c)
		return
	}
	utils.OkDetailed(response.PageResult{
		List:        list,
		Total:       total,
		PageSize:    pageInfo.PageSize,
		CurrentPage: pageInfo.CurrentPage,
	}, "获取成功", c)
}
