package model

import (
	"errors"
	"strings"
	"time"
	"yc_cron/model/request"
)

type Task struct {
	TaskId          int    `json:"task_id" gorm:"PRIMARY_KEY"`
	GroupId         int    `json:"group_id"`
	TaskName        string `json:"task_name"`
	Description     string `json:"description"`
	CronSpec        string `json:"cron_spec"`
	Command         string `json:"command"`
	Status          int    `json:"status"`
	LastExecuteTime int64  `json:"last_execute_time"`
	NextExecuteTime int64  `json:"next_execute_time"`
	CreateUserid    int    `json:"create_userid"`
	UpdateUserid    int    `json:"update_userid"`
	CreateTime      int64  `json:"create_time"`
	UpdateTime      int64  `json:"update_time"`
	Group           Group
}

const (
	TaskPrePare  = 1 //初始状态
	TaskStarting = 2 //已开始执行
)

func GetTaskList(search *request.BasePageInfo) (taskList []Task, total int, err error) {
	db := mdb
	if search.Condition != "" && search.SearchValue != "" {
		db = db.Where(search.Condition+" = ?", search.SearchValue)
	}
	if err = db.Model(&taskList).Count(&total).Error; err != nil {
		return
	}
	err = db.Limit(search.PageSize).Offset(search.PageSize * (search.CurrentPage - 1)).Find(&taskList).Error
	for idx, tl := range taskList {
		mdb.Where("group_id = ?", tl.GroupId).Find(&taskList[idx].Group)
	}
	return
}

func DelTask(taskId int) (err error) {
	var task Task
	if err = mdb.Delete(&task, taskId).Error; err != nil {
		return
	}
	return
}

func SaveTask(task *Task) (err error) {
	nowTime := time.Now().Unix()
	if task.TaskId == 0 {
		task.CreateTime = nowTime
		task.UpdateTime = nowTime
		err = mdb.Save(task).Error
	} else {
		err = mdb.Model(task).Updates(task).Error
	}
	if err != nil {
		if strings.Index(err.Error(), "uni_task_name") != -1 {
			err = errors.New("该分组名称已存在")
		}
	}
	return
}
