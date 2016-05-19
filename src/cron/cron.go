package util

import (
	"encoding/json"
	"github.com/Dataman-Cloud/omega-es/src/cache"
	"github.com/Dataman-Cloud/omega-es/src/es"
	"github.com/Dataman-Cloud/omega-es/src/model"
	//log "github.com/cihub/seelog"
	"github.com/sluu99/uuid"
	"time"
)

var mindex string

const (
	refreshtime = 3
	dietime     = 7
	polltime    = 10
)

func init() {
	mindex = uuid.Rand().Hex()
	//go RefreshHost()
	//go CleanDieHost()
	go PollAlarm()
}

func PollAlarm() {
	for {
		select {
		case <-time.After(polltime * time.Second):
			tnow := time.Now().Unix()
			//ahost, err := cache.GetActiveHost()
			//log.Debug("-----:", ahost, err)
			schedus, err := cache.GetAllSchedu()
			if err != nil {
				continue
			}
			infos, err := cache.GetAllAlarmInfo()
			if err != nil || len(infos) == 0 {
				continue
			}
			for k, v := range infos {
				var alarm model.LogAlarm
				//var alarm model.CronInfo
				if err = json.Unmarshal([]byte(v), &alarm); err != nil {
					continue
				}
				if stime, ok := schedus[k]; !ok || tnow-stime >= int64(alarm.Ival)*60 {
					es.JobExec([]byte(v))
				}
			}
		}
	}
}

func CleanDieHost() {
	for {
		select {
		case <-time.After(refreshtime * time.Second):
			hosts, err := cache.GetHost()
			if err == nil && len(hosts) > 0 {
				var diehost []string
				for k, v := range hosts {
					if time.Now().Unix()-v > dietime {
						diehost = append(diehost, k)
					}
				}
				cache.DelDieHost(diehost)
			}
		}
	}
}

func RefreshHost() {
	for {
		select {
		case <-time.After(refreshtime * time.Second):
			cache.RefreshHost(mindex)
		}
	}
}
