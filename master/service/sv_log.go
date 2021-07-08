package service

import (
	"cron_master/model"
	"cron_master/model/request"
)

func GetLogList(search *request.ComPageInfo) ([]model.TaskLog, int, error) {
	return model.GetLogList(search)
}
