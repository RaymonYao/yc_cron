package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientV3 "go.etcd.io/etcd/client/v3"
	"strconv"
	"time"
	"yc_cron/config"
)

type Etcd struct {
	Client *clientV3.Client
	Kv     clientV3.KV
	Lease  clientV3.Lease
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

// DeleteJob 删除任务
func (eClient *Etcd) DeleteJob(Id int) (oldJob *Job, err error) {
	var (
		jobKey    string
		delResp   *clientV3.DeleteResponse
		oldJobObj Job
	)

	//etcd中保存任务的key
	jobKey = config.GConfig.EtcdConfig.JobCronPrefix + strconv.Itoa(Id)

	//从etcd中删除它
	if delResp, err = eClient.Kv.Delete(context.TODO(), jobKey, clientV3.WithPrevKV()); err != nil {
		return
	}

	//返回被删除的任务信息
	if len(delResp.PrevKvs) != 0 {
		//解析一下旧值，返回它
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}
