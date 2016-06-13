package main

import (
	"github.com/Dataman-Cloud/omega-es/src/config"
	_ "github.com/Dataman-Cloud/omega-es/src/cron"
	. "github.com/Dataman-Cloud/omega-es/src/es"
	_ "github.com/Dataman-Cloud/omega-es/src/logger"
	"github.com/Dataman-Cloud/omega-es/src/util"
	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"net/http"
	"time"
)

func main() {
	initEcho()
}

func initEcho() {
	log.Info("echo framework starting...")
	addr := ":" + config.GetString("port")
	e := echo.New()

	e.Use(mw.Recover(), mw.Logger())

	//e.Use(CrossDomain)

	es := e.Group("/api/v3", auth)
	{
		es.Post("/es/index", SearchIndex)
		es.Post("/es/context", SearchContext)

		es.Post("/alarm", CreateLogAlarm)
		es.Put("/alarm", UpdateLogAlarm)
		es.Delete("/alarm/:id", DeleteLogAlarm)
		es.Patch("/alarm/:id", StopLogAlarm)

		es.Get("/alarm/:id", GetLogAlarm)
		es.Get("/alarm", GetAlarms)
		es.Get("/alarm/scheduler", GetAlarmHistory)
	}

	e.Get("/api/v3/es/download/index", ExportIndex)
	e.Get("/api/v3/es/download/context", ExportContext)
	e.Get("/api/v3/health/log", Health)

	log.Info("listening server address: ", addr)
	s := &http.Server{
		Addr:           addr,
		Handler:        e,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func CrossDomain(c *echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Request().Header.Set("Access-Control-Allow-Credentials", "true")
	c.Request().Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Request().Header.Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, X-Requested-By, If-Modified-Since, X-File-Name, Cache-Control, X-XSRFToken, Authorization")
	c.Request().Header.Set("Content-Type", "application/json")
	if c.Request().Method == "OPTIONS" {
		c.String(204, "")
	}
	return nil
}

func auth(c *echo.Context) error {
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
		return nil
	} else {
		return echo.NewHTTPError(http.StatusUnauthorized, time.Now().Format(time.RFC3339Nano)+" "+"validation failed...")
	}
}
