package model

import (
	"cron_master/model/request"
	"strconv"
)

type TaskLog struct {
	LogId      int    `json:"log_id" gorm:"PRIMARY_KEY"`
	TaskId     int    `json:"task_id"`
	OutPut     string `json:"out_put"`
	Error      string `json:"error"`
	Status     int    `json:"status"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
	CreateTime int64  `json:"create_time"`
	RunWorker  string `json:"run_worker"`
	Task       Task
}

func GetLogList(search *request.ComPageInfo) (logList []TaskLog, total int, err error) {
	db := mdb
	if search.Status != -2 {
		db = db.Where("status = ?", search.Status)
	}
	if search.Condition != "" && search.SearchValue != "" {
		if search.Condition == "task_name" {
			var taskList []*Task
			mdb.Where("task_name like ?", "%"+search.SearchValue+"%").Find(&taskList)
			var taskArr []string
			for _, v := range taskList {
				taskArr = append(taskArr, strconv.Itoa(v.TaskId))
			}
			db = db.Where("task_id in (?)", taskArr)
		} else {
			db = db.Where(search.Condition+" like ?", "%"+search.SearchValue+"%")
		}
	}
	if err = db.Model(&logList).Count(&total).Error; err != nil {
		return
	}
	if search.PageSize == 0 {
		err = db.Find(&logList).Error
	} else {
		err = db.Order("log_id desc").Limit(search.PageSize).Offset(search.PageSize * (search.CurrentPage - 1)).Find(&logList).Error
		if err != nil {
			return
		}
		for idx, tl := range logList {
			mdb.Where("task_id = ?", tl.TaskId).Find(&logList[idx].Task)
		}
	}
	return
}
