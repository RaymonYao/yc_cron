package service

import (
	"cron_master/model"
	"cron_master/model/request"
)

func GetTaskList(search *request.BasePageInfo) ([]model.Task, int, error) {
	return model.GetTaskList(search)
}

func SaveTask(u *model.Task) error {
	return model.SaveTask(u)
}

func DelTask(taskId int) error {
	return model.DelTask(taskId)
}
