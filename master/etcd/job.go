package etcd

// Job 存放在ETCD上的Job定时任务结构
type Job struct {
	Id       int    `json:"id"`       //任务ID
	Name     string `json:"name"`     //任务名
	Command  string `json:"command"`  //shell命令
	CronSpec string `json:"cronSpec"` //cron表达式
}
