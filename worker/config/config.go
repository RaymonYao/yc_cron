package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

var (
	GConfig *SysConfig
	cfgLock sync.Mutex
)

type SysConfig struct {
	Env        string      `json:"env"`
	Port       int         `json:"port"`
	EtcdConfig *EtcdConfig `json:"etcd"`
	StartTime  time.Time   `json:"start_time"`
}

type EtcdConfig struct {
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
	JobCronPrefix   string   `json:"jobCronPrefix"`
	JobKillPrefix   string   `json:"jobKillPrefix"`
	JobLockPrefix   string   `json:"jobKillPrefix"`
}

func InitConfig() {
	cfgLock.Lock()
	defer cfgLock.Unlock()
	if GConfig != nil {
		return
	}
	var (
		bytes     []byte
		sysConfig SysConfig
		err       error
	)
	if bytes, err = ioutil.ReadFile("./config/config.json"); err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(bytes, &sysConfig); err != nil {
		log.Fatal(err)
	}
	sysConfig.StartTime = time.Now()
	GConfig = &sysConfig
}
