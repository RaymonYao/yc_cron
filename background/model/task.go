package model

import (
	"errors"
	"strings"
	"time"
	"yc_cron/model/request"
)

type Task struct {
	TaskId       int    `json:"task_id" gorm:"PRIMARY_KEY"`
	TaskName     string `json:"task_name"`
	Description  string `json:"description"`
	Group        Group
	CronSpec     string    `json:"cron_spec"`
	Command      string    `json:"command"`
	Status       int       `json:"status"`
	CreateUserid int       `json:"create_userid"`
	UpdateUserid int       `json:"update_userid"`
	CreateAt     time.Time `json:"create_at"`
	UpdateAt     time.Time `json:"update_at"`
}

const (
	TaskPrePare  = 1 //初始状态
	TaskStarting = 2 //已开始执行
)

func GetTaskList(search *request.ComPageInfo) (taskList []Task, total int, err error) {
	db := mdb
	if search.Condition != "" && search.SearchValue != "" {
		db = db.Where(search.Condition+" = ?", search.SearchValue)
	}
	if search.Status != 0 {
		db = db.Where("status = ?", search.Status)
	}
	if err = db.Model(&taskList).Count(&total).Error; err != nil {
		return
	}
	err = db.Limit(search.PageSize).Offset(search.PageSize * (search.CurrentPage - 1)).Find(&taskList).Error
	return
}

func DelTask(taskId int) (err error) {
	var task Task
	if err = mdb.Where("task_id = ?", taskId).First(&task).Error; err != nil {
		return
	}
	if err = mdb.Delete(&task).Error; err != nil {
		return
	}
	return
}

func SaveTask(p *Task) (err error) {
	if p.TaskId == 0 {
		p.CreateAt = time.Now()
		p.UpdateAt = time.Now()
		err = mdb.Save(p).Error
	} else {
		err = mdb.Model(p).Updates(p).Error
	}
	if err != nil {
		if strings.Index(err.Error(), "uni_task_name") != -1 {
			err = errors.New("该分组名称已存在")
		}
	}
	return
}

func GetTaskById(taskId int, isInit bool) (task Task, err error) {
	err = mdb.Where("task_id = ?", taskId).First(&task).Error
	return
}
