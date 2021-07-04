package yc_worker

import (
	"cron_worker/etcd"
	"time"
)

var (
	JobChan chan *etcd.Job
	JobResultChan chan *etcd.JobExecuteResult
)

// InitScheduler 初始化调度
func InitScheduler() (err error) {
	//jobChan = make(chan *etcd.JobExecuteResult, 1000)
	go func() {
		//启动调度协程
		var (
			job           *etcd.Job
			scheduleAfter time.Duration
			scheduleTimer *time.Timer
			jobResult     *etcd.JobExecuteResult
		)

		//初始化一次(1次)
		//scheduleAfter = scheduler.TrySchedule()

		//调度的延迟定时器
		scheduleTimer = time.NewTimer(scheduleAfter)

		//定时任务common.Job
		for {
			select {
			case job = <-JobChan: //监听任务管道
				//对内存中维护的任务列表做增删改查
			case <-scheduleTimer.C:
				//最近的任务到期了
			case jobResult = <-JobResultChan:
				//监听任务执行结果
			}
			//调度一次任务
			//scheduleAfter = scheduler.TrySchedule()
			//重置调度间隔
			scheduleTimer.Reset(scheduleAfter)
		}
	}()
	return
}
