package util

import (
	"github.com/Dataman-Cloud/omega-es/src/config"
	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func init() {
	pool = initPool()
}

func initPool() *redis.Pool {
	err, host := config.GetStringMapString("redis", "host")
	if err != nil {
		log.Error(err)
	}
	err, port := config.GetStringMapString("redis", "port")
	if err != nil {
		log.Warn("can't find redis port default:6379")
		port = "6379"
	}
	return redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", host+":"+port)
		return c, err
	}, 20)
}

func Open() redis.Conn {
	if pool != nil {
		return pool.Get()
	}
	pool = initPool()
	return pool.Get()
}
