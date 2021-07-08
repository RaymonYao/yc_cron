package model

import (
	"cron_master/etcd"
	"cron_master/model/request"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
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
	PrevExecuteTime int64  `json:"prev_execute_time"`
	NextExecuteTime int64  `json:"next_execute_time" gorm:"-"` //gorm:"-" 设置忽略某字段
	CreateUserid    int    `json:"create_userid"`
	UpdateUserid    int    `json:"update_userid"`
	CreateTime      int64  `json:"create_time"`
	UpdateTime      int64  `json:"update_time"`
	Group           Group
}

func GetTaskList(search *request.ComPageInfo) (taskList []Task, total int, err error) {
	db := mdb
	if search.Status != -1 {
		db = db.Where("status = ?", search.Status)
	}
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
			//实时计算下次任务的执行时间
			p := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
			schedule, _ := p.Parse(tl.CronSpec)
			taskList[idx].NextExecuteTime = schedule.Next(time.Now()).Unix()
		}
	}
	return
}

func DelTask(taskId int) (err error) {
	db := mdb
	err = db.Transaction(func(tx *gorm.DB) (e error) {
		var task Task
		if e = tx.Delete(&task, taskId).Error; e != nil {
			return
		}
		_, e = etcd.EClient.DeleteJob(taskId)
		return
	})
	return
}

func PauseTask(taskId int) (err error) {
	db := mdb
	err = db.Transaction(func(tx *gorm.DB) (e error) {
		//gorm对于空值(0, nil, "", false)这些会被忽略掉，这里改用map来更新
		//data := make(map[string]interface{})
		//data["status"] = 0
		//data := map[string]interface{}{"status": 0}
		if e = tx.Model(Task{TaskId: taskId}).Updates(map[string]interface{}{"status": 0}).Error; e != nil {
			return
		}
		_, e = etcd.EClient.DeleteJob(taskId)
		return
	})
	return
}

func StartTask(taskId int) (err error) {
	db := mdb
	err = db.Transaction(func(tx *gorm.DB) (e error) {
		task := &Task{TaskId: taskId, Status: 1}
		if e = tx.Model(&task).Updates(&task).Error; e != nil {
			return
		}
		tx.Where("task_id = ?", taskId).Find(&task)
		_, err = etcd.EClient.SaveJob(&etcd.Job{
			Id:       task.TaskId,
			Name:     task.TaskName,
			Command:  task.Command,
			CronSpec: task.CronSpec,
		})
		return
	})
	return
}

func RunTask(taskId int) (err error) {
	db := mdb
	task := &Task{TaskId: taskId}
	db.Where("task_id = ?", taskId).Find(&task)
	err = etcd.EClient.RunJob(&etcd.Job{
		Id:       task.TaskId,
		Name:     task.TaskName,
		Command:  task.Command,
		CronSpec: task.CronSpec,
	})
	if err != nil {
		return
	}
	return
}

func KillTask(taskId int) (err error) {
	err = etcd.EClient.KillJob(taskId)
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
				e = errors.New("该任务名称已存在")
			}
		} else {
			if task.Status == 1 {
				_, e = etcd.EClient.SaveJob(&etcd.Job{
					Id:       task.TaskId,
					Name:     task.TaskName,
					Command:  task.Command,
					CronSpec: task.CronSpec,
				})
			}
		}
		return
	})
	return
}

func InitEtcdJob() {
	var (
		taskList []Task
	)
	taskList, _, _ = GetTaskList(&request.ComPageInfo{Status: 1})
	for _, task := range taskList {
		_, _ = etcd.EClient.SaveJob(&etcd.Job{
			Id:       task.TaskId,
			Name:     task.TaskName,
			Command:  task.Command,
			CronSpec: task.CronSpec,
		})
	}
	return
}
