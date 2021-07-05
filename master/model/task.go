package model

import (
	"cron_master/etcd"
	"cron_master/model/request"
	"errors"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
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
	if search.PageSize == 0 {
		err = db.Find(&taskList).Error
	} else {
		err = db.Limit(search.PageSize).Offset(search.PageSize * (search.CurrentPage - 1)).Find(&taskList).Error
		for idx, tl := range taskList {
			mdb.Where("group_id = ?", tl.GroupId).Find(&taskList[idx].Group)
		}
	}
	return
}

func DelTask(taskId int) (err error) {
	db := mdb
	err = db.Transaction(func(tx *gorm.DB) (e error) {
		var task Task
		if e = tx.Delete(&task, taskId).Error; err != nil {
			return
		}
		_, e = etcd.EClient.DeleteJob(taskId)
		return
	})
	return
}

func SaveTask(task *Task) (err error) {
	db := mdb
	err = db.Transaction(func(tx *gorm.DB) (e error) {
		nowTime := time.Now().Unix()
		task.UpdateTime = nowTime
		if task.TaskId == 0 {
			task.CreateTime = nowTime
			e = tx.Save(task).Error
		} else {
			e = tx.Model(task).Updates(task).Error
		}
		if e != nil {
			if strings.Index(e.Error(), "uni_task_name") != -1 {
				e = errors.New("该分组名称已存在")
			}
		} else {
			_, e = etcd.EClient.SaveJob(&etcd.Job{
				Id:       task.TaskId,
				Name:     task.TaskName,
				Command:  task.Command,
				CronSpec: task.CronSpec,
			})
		}
		return
	})
	return
}

func InitEtcdJob() (err error) {
	var (
		taskList []Task
	)
	taskList, _, err = GetTaskList(&request.BasePageInfo{})
	for _, task := range taskList {
		etcd.EClient.SaveJob(&etcd.Job{
			Id:       task.TaskId,
			Name:     task.TaskName,
			Command:  task.Command,
			CronSpec: task.CronSpec,
		})
	}
	return
}
