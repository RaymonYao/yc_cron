package model

import (
	"cron_master/etcd"
	"cron_master/model/request"
)

type Worker struct {
	WorkerID int    `json:"worker_id"`
	WorkerIP string `json:"worker_ip"`
}

func GetWorkerList(search *request.ComPageInfo) (workerList []Worker, total int, err error) {
	var workerListArr []string
	if workerListArr, err = etcd.EClient.ListWorkers(); err != nil {
		return
	}
	for k, v := range workerListArr {
		workerList = append(workerList, Worker{WorkerID: k + 1, WorkerIP: v})
	}
	return
}
