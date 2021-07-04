package service

import (
	"cron_master/model"
	"cron_master/model/request"
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
