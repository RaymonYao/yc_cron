package service

import (
	"yc_cron/model"
	"yc_cron/model/request"
)

func GetGroupList(search *request.BasePageInfo) ([]model.Group, int, error) {
	return model.GetGroupList(search)
}

func SaveGroup(u *model.Group) error {
	return model.SaveGroup(u)
}

func DelGroup(taskId int) error {
	return model.DelGroup(taskId)
}
