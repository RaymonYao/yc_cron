package api

import (
	"cron_master/model"
	"cron_master/model/request"
	"cron_master/model/response"
	"cron_master/service"
	"cron_master/utils"
	"github.com/gin-gonic/gin"
	"strconv"
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

func DelWorker(c *gin.Context) {
	var (
		workerId int
		err      error
	)
	if workerId, err = strconv.Atoi(c.Query("worker_id")); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	if err = service.DelWorker(workerId); err != nil {
		utils.FailWithMessage("删除失败, 原因:"+err.Error(), c)
		return
	}
	utils.OkWithMessage("删除成功", c)
}

func SaveWorker(c *gin.Context) {
	var (
		worker       model.Worker
		workerVerify map[string][]string
		err          error
	)
	_ = c.ShouldBindJSON(&worker)
	workerVerify = utils.Rules{
		"WorkerIP":    {utils.NotEmpty()},
		"Description": {utils.NotEmpty()},
	}
	if err = utils.Verify(worker, workerVerify); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	if err = service.SaveWorker(&worker); err != nil {
		utils.FailWithMessage("保存失败,"+err.Error(), c)
		return
	}
	utils.OkWithMessage("保存成功", c)
}
