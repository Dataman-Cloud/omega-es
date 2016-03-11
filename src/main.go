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
	addr := config.GetString("host") + ":" + config.GetString("port")
	e := echo.New()

	e.Use(mw.Recover(), mw.Logger())

	//es := e.Group("/es")
	es := e.Group("/es", auth)
	{
		es.Post("/index", SearchIndex)
		es.Post("/index/download", IndexExport)
		es.Post("/context", SearchContext)
		es.Post("/context/download", ContextExport)
		es.Post("/seniorsearch/:userId", SeniorSearch)
		es.Get("/appagg/:userId", AppAgg)
		es.Get("/topagg/:field/:userId", TopAgg)
	}
	download := e.Group("/es/download")
	{
		es.Get("/index/log", ExportIndex)
		es.Get("/context/log", ExportContext)
	}

	api := e.Group("/api/v3")
	{
		api.Get("/health/log", Health)
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
