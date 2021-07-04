package etcd

import (
	"context"
	yc_worker "cron_worker"
	"cron_worker/config"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientV3 "go.etcd.io/etcd/client/v3"
	"strconv"
	"time"
)

type Etcd struct {
	Client *clientV3.Client
	Kv     clientV3.KV
	Lease  clientV3.Lease
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

	EClient = &Etcd{
		Client: client,
		Kv:     kv,
		Lease:  lease,
	}
	EClient.WatchJobs()
	return
}

// SaveJob 保存任务
func (eClient *Etcd) SaveJob(job *Job) (oldJob *Job, err error) {
	var (
		jobKey    string
		jobValue  []byte
		putResp   *clientV3.PutResponse
		oldJobObj Job
	)

	//etcd的保存key
	jobKey = config.GConfig.EtcdConfig.JobCronPrefix + strconv.Itoa(job.Id)
	//任务信息json
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}
	//保存到etcd
	if putResp, err = eClient.Kv.Put(context.TODO(), jobKey, string(jobValue), clientV3.WithPrevKV()); err != nil {
		fmt.Println(err)
		return
	}
	//如果是更新，那么返回旧值
	if putResp.PrevKv != nil {
		//对旧值做一个反序列化
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

//监听任务变化
func (eClient *Etcd) WatchJobs() (err error) {
	var (
		getResp            *clientV3.GetResponse
		kvPair             *mvccpb.KeyValue
		job                *Job
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
		var job *Job
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			return
		}
		yc_worker.JobChan <- job
	}

	//2,从该revision向后监听变化事件
	go func() {
		//监听协程
		//从Get时刻的后续版本开始监听变化
		watchStartRevision = getResp.Header.Revision + 1
		//监听job_cron_前缀的任务后续变化
		watchChan = eClient.Watcher.Watch(context.TODO(), config.GConfig.EtcdConfig.JobCronPrefix, clientV3.WithRev(watchStartRevision), clientV3.WithPrefix())
		//处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					//任务保存事件
					if err = json.Unmarshal(watchEvent.Kv.Value, job); err != nil {
						continue
					}
					//构建一个更新Event
					//jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				case mvccpb.DELETE:

				}
				//变化推给scheduler
				//GScheduler.PushJobEvent(jobEvent)
			}
		}
	}()
	return
}
