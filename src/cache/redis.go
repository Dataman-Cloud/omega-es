package cache

import (
	"encoding/json"
	"errors"
	"github.com/Dataman-Cloud/omega-es/src/model"
	"github.com/Dataman-Cloud/omega-es/src/util"
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	alarmindex   = "alarm:set"
	alarminfo    = "alarm:info"
	alarmhost    = "alarm:host"
	alarmschedu  = "alarm:schedu"
	alarmscaling = "alarm:scaling"
)

func AddAlarm(id int64, info []byte) error {
	conn := util.Open()
	defer conn.Close()
	/*_, err := conn.Do("SADD", alarmindex, alarm.Id)
	if err != nil {
		return err
	}*/
	_, err := conn.Do("HSET", alarminfo, id, string(info))
	if err != nil {
		return err
	}
	return nil
}

func DeleteAlarm(id int64) error {
	conn := util.Open()
	defer conn.Close()
	if exists, err := redis.Bool(conn.Do("HEXISTS", alarminfo, id)); err != nil {
		return err
	} else if exists {
		if _, err := conn.Do("HDEL", alarminfo, id); err != nil {
			return err
		}
	}
	if exists, err := redis.Bool(conn.Do("HEXISTS", alarmschedu, id)); err != nil {
		return err
	} else if exists {
		if _, err := conn.Do("HDEL", alarmschedu, id); err != nil {
			return err
		}
	}
	/*if _, err := conn.Do("SREM", alarmindex, id); err != nil {
		return err
	}*/
	return nil
}

func UpdateAlarm(alarm *model.LogAlarm) error {
	conn := util.Open()
	defer conn.Close()
	if exists, err := redis.Bool(conn.Do("HEXISTS", alarminfo, alarm.Id)); err != nil {
		return err
	} else if !exists {
		/*if _, err = conn.Do("SADD", alarmindex, alarm.Id); err != nil {
			return err
		}*/
	}
	info, err := json.Marshal(alarm)
	if err != nil {
		return err
	}
	if _, err = conn.Do("HSET", alarminfo, alarm.Id, string(info)); err != nil {
		return err
	}
	return nil
}

func RefreshHost(host string) error {
	conn := util.Open()
	defer conn.Close()
	if _, err := conn.Do("HSET", alarmhost, host, time.Now().Unix()); err != nil {
		return err
	}
	return nil
}

func GetHost() (map[string]int64, error) {
	conn := util.Open()
	defer conn.Close()
	hosts, err := redis.Int64Map(conn.Do("HGETALL", alarmhost))
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

func DelDieHost(hosts []string) error {
	conn := util.Open()
	defer conn.Close()
	for _, v := range hosts {
		conn.Do("HDEL", alarmhost, v)
	}
	/*if _, err := conn.Do("HDEL", alarmhost, hosts...); err != nil {
		return err
	}*/
	return nil
}

func GetActiveHost() ([]string, error) {
	conn := util.Open()
	defer conn.Close()
	ahosts, err := redis.Strings(conn.Do("HKEYS", alarmhost))
	return ahosts, err
}

func GetAllAlarmInfo() (map[string]string, error) {
	conn := util.Open()
	defer conn.Close()
	infos, err := redis.StringMap(conn.Do("HGETALL", alarminfo))
	return infos, err
}

func GetAllSchedu() (map[string]int64, error) {
	conn := util.Open()
	defer conn.Close()
	return redis.Int64Map(conn.Do("HGETALL", alarmschedu))
}

func UpdateScheduTime(id int64) {
	conn := util.Open()
	defer conn.Close()
	conn.Do("HSET", alarmschedu, id, time.Now().Unix())
}

func AppExtend(appid int64, maxs uint64) error {
	conn := util.Open()
	defer conn.Close()
	m, err := redis.Uint64(conn.Do("HGET", alarmscaling, appid))
	if err == nil {
		return err
	}
	if m == maxs {
		return errors.New("Have reached the maximum")
	}
	_, err = conn.Do("HSET", alarmscaling, appid, maxs)
	return err
}

func AppShrink(appid int64, min64 uint64) error {
	conn := util.Open()
	defer conn.Close()
	m, err := redis.Uint64(conn.Do("HGET", alarmscaling, appid))
	if err != nil {
		return err
	}
	if m <= min64 {
		return errors.New("Have reached a minimum value")
	}
	_, err = conn.Do("HSET", alarmscaling, appid, m-1)
	return err
}
