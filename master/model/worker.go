package model

import (
	"cron_master/etcd"
	"cron_master/model/request"
	"errors"
	"strings"
	"time"
)

type Worker struct {
	WorkerId     uint   `json:"worker_id" gorm:"PRIMARY_KEY"`
	WorkerIP     string `json:"worker_ip"`
	Description  string `json:"description"`
	Status       int    `json:"status" gorm:"-"` //gorm:"-" 设置忽略某字段
	CreateUserid int    `json:"create_userid"`
	UpdateUserid int    `json:"update_userid"`
	CreateTime   int64  `json:"create_time"`
	UpdateTime   int64  `json:"update_time"`
}

func GetWorkerList(search *request.ComPageInfo) (workerList []Worker, total int, err error) {
	db := mdb
	if search.Condition != "" && search.SearchValue != "" {
		db = db.Where(search.Condition+" like ?", "%"+search.SearchValue+"%")
	}
	if err = db.Model(&workerList).Count(&total).Error; err != nil {
		return
	}
	err = db.Order("worker_id desc").Limit(search.PageSize).Offset(search.PageSize * (search.CurrentPage - 1)).Find(&workerList).Error
	if err != nil {
		return
	}
	var workerListArr []string
	workerListArr, err = etcd.EClient.ListWorkers()
	for idx, wl := range workerList {
		if len(workerListArr) == 0 {
			workerList[idx].Status = -1
		} else {
			for _, v := range workerListArr {
				if v == wl.WorkerIP {
					workerList[idx].Status = 0
					break
				} else {
					workerList[idx].Status = -1
				}
			}
		}
	}
	return
}

func DelWorker(workerId int) (err error) {
	db := mdb
	var worker Worker
	if err = db.Delete(&worker, workerId).Error; err != nil {
		return
	}
	return
}

func SaveWorker(worker *Worker) (err error) {
	nowTime := time.Now().Unix()
	if worker.WorkerId == 0 {
		worker.CreateTime = nowTime
		worker.UpdateTime = nowTime
		err = mdb.Save(worker).Error
	} else {
		err = mdb.Model(worker).Updates(worker).Error
	}
	if err != nil {
		if strings.Index(err.Error(), "uni_worker_ip") != -1 {
			err = errors.New("该节点IP已存在")
		}
	}
	return
}
