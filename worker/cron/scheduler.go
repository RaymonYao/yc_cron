package cron

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"os/exec"
	"time"
)

var (
	job              *Job
	JobScheduleTable map[int]*Job
	JobResultChan    chan *JobExecuteResult
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
				delete(JobExecutingTable, jobResult.Job.Id)
				//todo写日志
			}
		}
	}()
	return
}

// AddCron 调度队列中添加任务
func AddCron(job *Job) {
	if JobScheduleTable[job.Id] != nil {
		oldJob := JobScheduleTable[job.Id]
		fmt.Println(oldJob.EntryId)
		GCron.Remove(oldJob.EntryId)
		delete(JobScheduleTable, oldJob.Id)
	}

	EntryId, err = GCron.AddFunc(job.CronSpec, func() {
		ExecuteJob(job)
		cmd := exec.CommandContext(context.TODO(), "/bin/bash", "-c", job.Command)
		output, _ := cmd.CombinedOutput()
		fmt.Println(string(output))
		fmt.Println(time.Now())
		fmt.Println("\n")
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
