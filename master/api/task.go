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

func GetTaskList(c *gin.Context) {
	var (
		pageInfo request.ComPageInfo
		list     []model.Task
		total    int
		err      error
	)
	_ = c.ShouldBindJSON(&pageInfo)
	if list, total, err = service.GetTaskList(&pageInfo); err != nil {
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

func DelTask(c *gin.Context) {
	var (
		taskId int
		err    error
	)
	if taskId, err = strconv.Atoi(c.Query("task_id")); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	if err = service.DelTask(taskId); err != nil {
		utils.FailWithMessage("删除失败, 原因:"+err.Error(), c)
		return
	}
	utils.OkWithMessage("删除成功", c)
}

func PauseTask(c *gin.Context) {
	var (
		taskId int
		err    error
	)
	if taskId, err = strconv.Atoi(c.Query("task_id")); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	if err = service.PauseTask(taskId); err != nil {
		utils.FailWithMessage("暂停失败, 原因:"+err.Error(), c)
		return
	}
	utils.OkWithMessage("暂停成功", c)
}

func StartTask(c *gin.Context) {
	var (
		taskId int
		err    error
	)
	if taskId, err = strconv.Atoi(c.Query("task_id")); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	if err = service.StartTask(taskId); err != nil {
		utils.FailWithMessage("开启失败, 原因:"+err.Error(), c)
		return
	}
	utils.OkWithMessage("开启成功", c)
}

func RunTask(c *gin.Context) {
	var (
		taskId int
		err    error
	)
	if taskId, err = strconv.Atoi(c.Query("task_id")); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	if err = service.RunTask(taskId); err != nil {
		utils.FailWithMessage("执行失败, 原因:"+err.Error(), c)
		return
	}
	utils.OkWithMessage("已在后台执行", c)
}

func KillTask(c *gin.Context) {
	var (
		taskId int
		err    error
	)
	if taskId, err = strconv.Atoi(c.Query("task_id")); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	if err = service.KillTask(taskId); err != nil {
		utils.FailWithMessage("执行失败, 原因:"+err.Error(), c)
		return
	}
	utils.OkWithMessage("已在后台强杀", c)
}

func SaveTask(c *gin.Context) {
	var (
		task       model.Task
		taskVerify map[string][]string
		err        error
	)
	_ = c.ShouldBindJSON(&task)
	taskVerify = utils.Rules{
		"TaskName":    {utils.NotEmpty(), utils.Le("200"), utils.Ge("3")},
		"Description": {utils.NotEmpty()},
		"CronSpec":    {utils.NotEmpty()},
		"Command":     {utils.NotEmpty()},
	}
	if err = utils.Verify(task, taskVerify); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	if err = service.SaveTask(&task); err != nil {
		utils.FailWithMessage("保存失败,"+err.Error(), c)
		return
	}
	utils.OkWithMessage("保存成功", c)
}
