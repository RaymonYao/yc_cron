package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"yc_cron/model"
	"yc_cron/model/request"
	"yc_cron/model/response"
	"yc_cron/service"
	"yc_cron/utils"
)

func GetGroupList(c *gin.Context) {
	var (
		pageInfo request.BasePageInfo
		list     []model.Group
		total    int
		err      error
	)
	_ = c.ShouldBindJSON(&pageInfo)
	if list, total, err = service.GetGroupList(&pageInfo); err != nil {
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

func DelGroup(c *gin.Context) {
	var (
		groupId int
		err     error
	)
	if groupId, err = strconv.Atoi(c.Query("group_id")); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	if err = service.DelGroup(groupId); err != nil {
		utils.FailWithMessage("删除失败, 原因:"+err.Error(), c)
		return
	}
	utils.OkWithMessage("删除成功", c)
}

func SaveGroup(c *gin.Context) {
	var (
		group       model.Group
		groupVerify map[string][]string
		err         error
	)
	_ = c.ShouldBindJSON(&group)
	groupVerify = utils.Rules{
		"UserName": {utils.NotEmpty(), utils.Le("20")},
		"Password": {utils.Le("18")},
		"NickName": {utils.NotEmpty(), utils.Le("20")},
	}
	if err = utils.Verify(group, groupVerify); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	if err = service.SaveGroup(&group); err != nil {
		utils.FailWithMessage("保存失败,"+err.Error(), c)
		return
	}
	utils.OkWithMessage("保存成功", c)
}
