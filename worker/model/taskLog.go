package model

import "time"

type TaskLog struct {
	LogId       int    `json:"log_id" gorm:"PRIMARY_KEY"`
	TaskId      int    `json:"task_id"`
	OutPut      string `json:"out_put"`
	Error       string `json:"error"`
	Status      int    `json:"status"`
	ProcessTime int64  `json:"process_time"`
	CreateTime  int64  `json:"create_time"`
}

func SaveTaskLog(taskLog *TaskLog) (err error) {
	db := mdb
	taskLog.CreateTime = time.Now().Unix()
	err = db.Save(taskLog).Error
	return
}
