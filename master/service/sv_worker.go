package service

import (
	"cron_master/model"
	"cron_master/model/request"
)

func GetWorkerList(search *request.ComPageInfo) ([]model.Worker, int, error) {
	return model.GetWorkerList(search)
}
