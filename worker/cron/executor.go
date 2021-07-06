package cron

import (
	"context"
	"log"
	"math/rand"
	"os/exec"
	"time"
)

var (
	JobExecutingTable map[int]*Job
	jobResultChan     chan *JobExecuteResult //任务结果管道
)

// InitExecutor 初始化执行器
func InitExecutor() {
	JobExecutingTable = make(map[int]*Job)
	JobResultChan = make(chan *JobExecuteResult, 1000)
}

// ExecuteJob 执行任务
func ExecuteJob(job *Job) {
	//如果任务正在执行，跳过本次调度(执行的任务可能运行很久，1分钟会调度很多次，但是只能执行1次，防止并发!)
	if JobExecutingTable[job.Id] != nil {
		log.Println("任务还未执行完:", job.Name)
		return
	}

	JobExecutingTable[job.Id] = job

	//开启协程执行任务
	go func() {
		var (
			cmd     *exec.Cmd
			err     error
			output  []byte
			result  *JobExecuteResult
			jobLock *JobLock
		)

		//任务结果
		result = &JobExecuteResult{
			Job:    job,
			Output: make([]byte, 0),
		}

		//初始化分布式锁
		jobLock = &JobLock{Job: job}

		//记录任务开始时间
		result.StartTime = time.Now()

		//上锁
		//随机睡眠(0-1s) 为了防止各个服务器时间不准导致的抢锁不公平，正常情况下各个服务器会用ntp时间服务器进行时间同步
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		jobLock, err = TryLock(job)
		defer jobLock.Unlock()

		if err != nil {
			//上锁失败
			result.Err = err
			result.EndTime = time.Now()
		} else {
			//上锁成功后，重置任务启动时间
			result.StartTime = time.Now()

			//用于取消任务执行
			cancelCtx, cancelFunc := context.WithCancel(context.TODO())

			//执行shell命令
			cmd = exec.CommandContext(cancelCtx, "/bin/bash", "-c", job.Command)

			//执行并捕获输出
			output, err = cmd.CombinedOutput()

			//记录任务结束时间
			result.Job = job
			result.EndTime = time.Now()
			result.Output = output
			result.Err = err
			job.CancelFunc = cancelFunc
		}
		//任务执行完成后，把执行的结果返回给Scheduler,Scheduler会从executingTable中删除掉执行记录
		jobResultChan <- result
	}()
	return
}
