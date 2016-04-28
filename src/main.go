package main

import (
	"github.com/Dataman-Cloud/omega-es/src/config"
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

	es := e.Group("/es", auth)
	{
		es.Post("/index", SearchIndex)
		es.Post("/context", SearchContext)
		es.Post("/seniorsearch/:userId", SeniorSearch)
		es.Get("/appagg/:userId", AppAgg)
		es.Get("/topagg/:field/:userId", TopAgg)
	}

	ea := e.Group("/es/alarm", auth)
	{
		ea.Post("/create", CreateLogAlarm)
		ea.Delete("/delete/:id", DeleteAlarm)
		ea.Get("/list", GetAlarms)
		ea.Get("/scheduler/history", GetAlarmHistory)
	}

	/*ew := e.Group("/es/watcher")
	{
		ew.Post("/create", CreateWatcher)
		ew.Post("/delete", DeleteWatcher)
		ew.Get("/list/:usertype/:userid", GetWatchers)
		ew.Post("/history", GetWatcherHistory)
	}*/

	ed := e.Group("/es/download")
	{
		ed.Get("/index/log", ExportIndex)
		ed.Get("/context/log", ExportContext)
	}

	api := e.Group("/api/v3")
	{
		api.Get("/health/log", Health)
		api.Post("/scheduler", JobExec)
	}

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
