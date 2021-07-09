package model

import (
	"fmt"
	"time"
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
}

func SaveTaskLog(taskLog *TaskLog) (err error) {
	db := mdb
	taskLog.CreateTime = time.Now().Unix()
	fmt.Println(taskLog)
	err = db.Save(taskLog).Error
	return
}
