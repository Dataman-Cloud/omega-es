package util

import (
	"fmt"
	"github.com/Dataman-Cloud/omega-es/src/config"
	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func init() {
	pool = initPool()
}

func initPool() *redis.Pool {
	return redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp",
			fmt.Sprintf("%s:%d", config.GetConfig().Rc.Host,
				config.GetConfig().Rc.Port))
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
