package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dataman-Cloud/omega-es/src/config"
	_ "github.com/Dataman-Cloud/omega-es/src/cron"
	. "github.com/Dataman-Cloud/omega-es/src/es"
	_ "github.com/Dataman-Cloud/omega-es/src/logger"
	"github.com/Dataman-Cloud/omega-es/src/util"
	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConfigInit()
	util.ChronosInit()
	util.MysqlInit()
	util.EsInit()
	util.RedisInit()
	initEcho()
}

func initEcho() {
	log.Info("echo framework starting...")
	addr := fmt.Sprintf(":%d", config.GetConfig().Port)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.Logger(), gin.Recovery())
	v3 := r.Group("/api/v3", auth)
	{
		v3.POST("/es/index", SearchIndex)
		v3.POST("/es/context", SearchContext)

		v3.POST("/alarm", CreateLogAlarm)
		v3.PUT("/alarm", UpdateLogAlarm)
		v3.GET("/alarm", GetAlarms)

		v3.GET("/alarm/:id", GetLogAlarm)
		v3.PATCH("/alarm/:id", StopLogAlarm)
		v3.DELETE("/alarm/:id", DeleteLogAlarm)

		v3.GET("/alarms", GetAlarmHistory)
	}

	r.GET("/api/v3/health/log", Health)
	r.GET("/api/v3/es/download/index", ExportIndex)
	r.GET("/api/v3/es/download/context", ExportContext)

	log.Info("listening server address: ", addr)
	s := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func auth(c *gin.Context) {
	auth := false
	if authtoken := util.Header(c, "Authorization"); authtoken != "" {
		conn := util.Open()
		defer conn.Close()
		uid, err := redis.String(conn.Do("HGET", "s:"+authtoken, "user_id"))
		if err == nil {
			auth = true
			c.Set("uid", uid)
		} else if err != redis.ErrNil {
			log.Error("[app] got error1 ", err)
		}
	}
	if auth {
		c.Next()
	} else {
		c.String(http.StatusUnauthorized, "Invalid authentication")
		c.Abort()
	}
}
