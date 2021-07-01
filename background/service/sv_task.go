package service

import (
	"yc_cron/model"
	"yc_cron/model/request"
)

func GetTaskList(search *request.ComPageInfo) ([]model.Task, int, error) {
	return model.GetTaskList(search)
}

func SaveTask(u *model.Task) error {
	return model.SaveTask(u)
}

func DelTask(taskId int) error {
	return model.DelTask(taskId)
}
