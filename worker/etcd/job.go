package etcd

import (
	"context"
	"time"
)

// Job 存放在ETCD上的Job定时任务结构
type Job struct {
	Id       int    `json:"id"`       //任务ID
	Name     string `json:"name"`     //任务名
	Command  string `json:"command"`  //shell命令
	CronSpec string `json:"cronSpec"` //cron表达式
}

// JobExecuteInfo 任务执行状态
type JobExecuteInfo struct {
	Job        *Job               //任务信息
	PlanTime   time.Time          //理论上的调度时间
	RealTime   time.Time          //实际的调度时间
	CancelCtx  context.Context    //任务command的context
	CancelFunc context.CancelFunc //用于取消command执行的cancel函数
}

// JobExecuteResult 任务执行结果
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo //执行状态
	Output      []byte          //脚本输出
	Err         error           //脚本错误原因
	StartTime   time.Time       //启动时间
	EndTime     time.Time       //结束时间
}
