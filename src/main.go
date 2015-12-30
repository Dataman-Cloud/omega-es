package main

import (
	"errors"
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

	//e.Use(mw.Recover(), mw.Logger())
	e.Use(mw.Recover(), mw.Logger(), auth)

	es := e.Group("/es")
	{
		es.Post("/index", SearchIndex)
		es.Post("/context", SearchContext)
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
		err := errors.New("validation failed...")
		return err
	}
}
