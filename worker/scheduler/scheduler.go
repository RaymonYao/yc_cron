package scheduler

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"os/exec"
	"time"
)

var (
	job           *Job
	JobTable      map[int]*Job
	JobResultChan chan *JobExecuteResult
	err           error
	GCron         *cron.Cron
	EntryId       cron.EntryID
	jobExisted    bool
)

// InitScheduler 初始化调度
func InitScheduler() {
	JobTable = make(map[int]*Job)
	GCron = cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)))
	go func() {
		var (
			jobResult *JobExecuteResult
		)
		for {
			select {
			case jobResult = <-JobResultChan:
				//监听任务执行结果
				fmt.Println(jobResult)
			}
		}
	}()
	return
}

// AddCron 调度队列中添加任务
func AddCron(job *Job) {
	if JobTable[job.Id] != nil {
		oldJob := JobTable[job.Id]
		fmt.Println(oldJob.EntryId)
		GCron.Remove(oldJob.EntryId)
		delete(JobTable, oldJob.Id)
	}

	EntryId, err = GCron.AddFunc(job.CronSpec, func() {
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
	JobTable[job.Id] = job
	return
}

// RemoveCron 调度队列中移除任务
func RemoveCron(jobId int) {
	if JobTable[jobId] != nil {
		job = JobTable[jobId]
		GCron.Remove(job.EntryId)
		delete(JobTable, jobId)
	}
	return
}
