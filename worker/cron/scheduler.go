package cron

import (
	"cron_worker/model"
	"fmt"
	"github.com/robfig/cron/v3"
)

var (
	job              *Job
	JobScheduleTable map[int]*Job
	err              error
	GCron            *cron.Cron
	EntryId          cron.EntryID
)

// InitScheduler 初始化调度
func InitScheduler() {
	JobScheduleTable = make(map[int]*Job)
	GCron = cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)))
	go func() {
		for {
			select {
			//监听任务执行结果
			case jobResult := <-JobResultChan:
				fmt.Println(jobResult.Output)
				delete(JobExecutingTable, jobResult.Job.Id)
				_ = model.SaveTask(&model.Task{
					TaskId:          jobResult.Job.Id,
					PrevExecuteTime: jobResult.StartTime,
				})
				//写日志，暂时写进mysql，方便列表展示，后期会解耦，做成文件日志，用ELK去收集日志
				var (
					status int
					errStr string
				)
				if jobResult.Err != nil {
					status = -1
					errStr = jobResult.Err.Error()
				} else {
					status = 0
					errStr = ""
				}
				_ = model.SaveTaskLog(&model.TaskLog{
					TaskId:      jobResult.Job.Id,
					OutPut:      jobResult.Output,
					Error:       errStr,
					Status:      status,
					ProcessTime: jobResult.EndTime - jobResult.StartTime,
				})
			}
		}
	}()
	return
}

// AddCron 调度队列中添加任务
func AddCron(job *Job) {
	if JobScheduleTable[job.Id] != nil {
		oldJob := JobScheduleTable[job.Id]
		GCron.Remove(oldJob.EntryId)
		delete(JobScheduleTable, oldJob.Id)
	}

	EntryId, err = GCron.AddFunc(job.CronSpec, func() {
		ExecuteJob(job)
	})
	if err != nil {
		return
	}
	job.EntryId = EntryId
	JobScheduleTable[job.Id] = job
	return
}

// RemoveCron 调度队列中移除任务
func RemoveCron(jobId int) {
	if JobScheduleTable[jobId] != nil {
		job = JobScheduleTable[jobId]
		GCron.Remove(job.EntryId)
		delete(JobScheduleTable, jobId)
	}
	return
}
