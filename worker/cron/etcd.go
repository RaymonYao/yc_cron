package cron

import (
	"context"
	"cron_worker/config"
	"encoding/json"
	"errors"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientV3 "go.etcd.io/etcd/client/v3"
	"strconv"
	"strings"
	"time"
)

type Etcd struct {
	Client  *clientV3.Client
	Kv      clientV3.KV
	Lease   clientV3.Lease
	Watcher clientV3.Watcher
}

var (
	EClient *Etcd
)

func InitEtcd() {
	var (
		v3Config clientV3.Config
		client   *clientV3.Client
		kv       clientV3.KV
		lease    clientV3.Lease
		watcher  clientV3.Watcher
		err      error
	)

	//初始化配置
	v3Config = clientV3.Config{
		Endpoints:   config.GConfig.EtcdConfig.EtcdEndpoints,                                     //集群地址
		DialTimeout: time.Duration(config.GConfig.EtcdConfig.EtcdDialTimeout) * time.Millisecond, //连接超时
	}

	//建立连接
	if client, err = clientV3.New(v3Config); err != nil {
		return
	}

	//得到KV和Lease的API子集
	kv = clientV3.NewKV(client)
	lease = clientV3.NewLease(client)
	watcher = clientV3.NewWatcher(client)

	EClient = &Etcd{
		Client:  client,
		Kv:      kv,
		Lease:   lease,
		Watcher: watcher,
	}
	if err = EClient.WatchJobs(); err != nil {
		return
	}
	if err = EClient.WatchRunJobs(); err != nil {
		return
	}
	if err = EClient.WatchKillJobs(); err != nil {
		return
	}
	return
}

// WatchJobs 监听任务变化
func (eClient *Etcd) WatchJobs() (err error) {
	var (
		getResp            *clientV3.GetResponse
		kvPair             *mvccpb.KeyValue
		watchStartRevision int64
		watchChan          clientV3.WatchChan
		watchResp          clientV3.WatchResponse
		watchEvent         *clientV3.Event
	)

	//1, get一下job_cron_前缀的所有任务，并且获知当前集群的revision
	if getResp, err = eClient.Kv.Get(context.TODO(), config.GConfig.EtcdConfig.JobCronPrefix, clientV3.WithPrefix()); err != nil {
		return
	}
	//当前有哪些任务
	for _, kvPair = range getResp.Kvs {
		var job Job
		if err = json.Unmarshal(kvPair.Value, &job); err != nil {
			return
		}
		AddCron(&job)
	}
	GCron.Start()

	//2,从该revision向后监听变化事件
	go func() {
		for {
			//监听协程
			//从Get时刻的后续版本开始监听变化
			watchStartRevision = getResp.Header.Revision + 1
			//监听job_cron_前缀的任务后续变化
			watchChan = eClient.Watcher.Watch(context.TODO(), config.GConfig.EtcdConfig.JobCronPrefix, clientV3.WithRev(watchStartRevision), clientV3.WithPrefix())
			//处理监听事件
			for watchResp = range watchChan {
				for _, watchEvent = range watchResp.Events {
					var job Job
					//任务保存事件
					switch watchEvent.Type {
					case mvccpb.PUT:
						if err = json.Unmarshal(watchEvent.Kv.Value, &job); err != nil {
							continue
						}
						AddCron(&job)
					case mvccpb.DELETE:
						jobId, _ := strconv.Atoi(strings.TrimPrefix(string(watchEvent.Kv.Key), config.GConfig.EtcdConfig.JobCronPrefix))
						RemoveCron(jobId)
					}
				}
			}
		}
	}()
	return
}

// WatchRunJobs 监听手动执行的任务
func (eClient *Etcd) WatchRunJobs() (err error) {
	var (
		watchChan  clientV3.WatchChan
		watchResp  clientV3.WatchResponse
		watchEvent *clientV3.Event
	)

	go func() {
		for {
			watchChan = EClient.Watcher.Watch(context.TODO(), config.GConfig.EtcdConfig.JobRunPrefix, clientV3.WithPrefix())
			//处理监听事件
			for watchResp = range watchChan {
				for _, watchEvent = range watchResp.Events {
					var job Job
					//任务保存事件
					switch watchEvent.Type {
					case mvccpb.PUT:
						if err = json.Unmarshal(watchEvent.Kv.Value, &job); err != nil {
							continue
						}
						ExecuteJob(&job) //手动执行一次任务
					case mvccpb.DELETE:
						////job_run_租约过期，被自动删除
					}
				}
			}
		}
	}()
	return
}

// WatchKillJobs 监听被强杀的任务
func (eClient *Etcd) WatchKillJobs() (err error) {
	var (
		watchChan  clientV3.WatchChan
		watchResp  clientV3.WatchResponse
		watchEvent *clientV3.Event
	)

	//2,从该revision向后监听变化事件
	go func() {
		for {
			watchChan = EClient.Watcher.Watch(context.TODO(), config.GConfig.EtcdConfig.JobKillPrefix, clientV3.WithPrefix())
			//处理监听事件
			for watchResp = range watchChan {
				for _, watchEvent = range watchResp.Events {
					//强杀保存事件
					switch watchEvent.Type {
					case mvccpb.PUT:
						jobId, _ := strconv.Atoi(strings.TrimPrefix(string(watchEvent.Kv.Key), config.GConfig.EtcdConfig.JobKillPrefix))
						if JobExecutingTable[jobId] != nil {
							JobExecutingTable[jobId].CancelFunc() //执行强杀函数
						}
					case mvccpb.DELETE:
						////job_kill_租约过期，被自动删除
					}
				}
			}
		}
	}()
	return
}

func (jobLock *JobLock) TryLock(job *Job) (err error) {
	var (
		leaseGrantResp *clientV3.LeaseGrantResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		leaseId        clientV3.LeaseID
		keepRespChan   <-chan *clientV3.LeaseKeepAliveResponse
		txn            clientV3.Txn
		lockKey        string
		txnResp        *clientV3.TxnResponse
	)

	//1, 创建租约(5秒)
	if leaseGrantResp, err = EClient.Lease.Grant(context.TODO(), 5); err != nil {
		return
	}

	//context用于取消自动续租
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())

	//续租ID
	leaseId = leaseGrantResp.ID

	//2, 自动续租
	if keepRespChan, err = EClient.Lease.KeepAlive(cancelCtx, leaseId); err != nil {
		goto FAIL
	}

	//3, 处理续租应答的协程
	go func() {
		var (
			keepResp *clientV3.LeaseKeepAliveResponse
		)
		for {
			select {
			case keepResp = <-keepRespChan: //自动续租的应答
				if keepResp == nil {
					goto END
				}
			}
		}
	END:
	}()

	//4, 创建事务txn
	txn = EClient.Kv.Txn(context.TODO())

	//锁路径
	lockKey = config.GConfig.EtcdConfig.JobLockPrefix + strconv.Itoa(job.Id)

	//5,事务抢锁
	txn.If(clientV3.Compare(clientV3.CreateRevision(lockKey), "=", 0)).
		Then(clientV3.OpPut(lockKey, "", clientV3.WithLease(leaseId))).
		Else(clientV3.OpGet(lockKey))

	//提交事务
	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	//6,成功返回，失败释放租约
	if !txnResp.Succeeded {
		//锁被占用
		err = errors.New("锁已经被占用")
	}

	//抢锁成功
	jobLock.LeaseId = leaseId
	jobLock.CancelFunc = cancelFunc
	jobLock.IsLocked = true
	return

FAIL:
	cancelFunc()                                  //取消自动续租
	EClient.Lease.Revoke(context.TODO(), leaseId) //释放租约
	return
}

// Unlock 释放锁
func (jobLock *JobLock) Unlock() {
	if jobLock.IsLocked {
		jobLock.CancelFunc()                                  //取消我们程序自动续租的协程
		EClient.Lease.Revoke(context.TODO(), jobLock.LeaseId) //释放租约
	}
}
