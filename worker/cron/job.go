package cron

import (
	"context"
	"github.com/robfig/cron/v3"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// Job 存放在ETCD上的Job定时任务结构
type Job struct {
	Id         int    `json:"id"`       //任务ID
	Name       string `json:"name"`     //任务名
	Command    string `json:"command"`  //shell命令
	CronSpec   string `json:"cronSpec"` //cron表达式
	EntryId    cron.EntryID
	CancelFunc context.CancelFunc //用于取消command执行的cancel函数
}

// JobExecuteResult 任务执行结果
type JobExecuteResult struct {
	Job       *Job
	Output    []byte    //脚本输出
	Err       error     //脚本错误原因
	StartTime time.Time //启动时间
	EndTime   time.Time //结束时间
}

// JobLock 分布式锁(Txn事务)
type JobLock struct {
	Job        *Job
	CancelFunc context.CancelFunc //用于终止自动续租
	LeaseId    clientv3.LeaseID   //租约ID
	IsLocked   bool               //是否上锁成功
}
