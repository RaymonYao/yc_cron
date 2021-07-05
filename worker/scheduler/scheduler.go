package scheduler

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

var (
	job           *Job
	JobEventChan  chan string
	JobTable      map[int]*Job
	JobResultChan chan *JobExecuteResult
	err           error
	expr          *cronexpr.Expression
)

// InitScheduler 初始化调度
func InitScheduler() {
	JobTable = make(map[int]*Job)
	JobEventChan = make(chan string, 1000)
	go func() {
		//启动调度协程
		var (
			scheduleAfter time.Duration
			scheduleTimer *time.Timer
			jobResult     *JobExecuteResult
		)

		//调度的延迟定时器，设置一个初始时间，尽量设置大点，保证初始化任务队列能先调度一次，然后定时的时间会自动重置为最近一个要执行的任务离现在的时间间隔
		scheduleTimer = time.NewTimer(1000 * time.Second)

		for {
			select {
			case jobEventChan := <-JobEventChan: //监听任务管道
				//内存中维护的任务列表有变化
				fmt.Println(jobEventChan)
				for k, v := range JobTable {
					fmt.Println(k)
					fmt.Println(*v)
				}
			case <-scheduleTimer.C:
				//最近的任务到期了
			case jobResult = <-JobResultChan:
				//监听任务执行结果
				fmt.Println(jobResult)
			}
			fmt.Println("时间到")
			//调度一次任务
			scheduleAfter = TrySchedule()
			fmt.Println(scheduleAfter)
			//重置调度间隔
			scheduleTimer.Reset(10 * time.Second)
		}
	}()
	return
}

// TrySchedule 重新计算任务调度状态
func TrySchedule() (scheduleAfter time.Duration) {
	var (
		now      time.Time
		nextTime time.Time
		nearTime time.Time
	)

	//当前时间
	now = time.Now()

	//遍历所有任务
	for _, job = range JobTable {
		//解析Job的cron表达式
		if expr, err = cronexpr.Parse(job.CronSpec); err != nil {
			return
		}
		fmt.Println(expr)
		nextTime = expr.Next(now)
		fmt.Println(nextTime)
		if nextTime.Before(now) || nextTime.Equal(now) {
			//fmt.Println(nextTime)
			//TryStartJob(job)
			nextTime = expr.Next(time.Now()) //计算该任务下次要执行的时间
		}

		//统计最近一个要过期的任务时间
		if &nearTime == nil || nextTime.Before(nearTime) {
			nearTime = nextTime
		}
	}

	//下次调度间隔(最近要执行的任务调度时间-当前时间）
	scheduleAfter = (nearTime).Sub(now)
	return
}
