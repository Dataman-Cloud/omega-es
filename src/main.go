package main

import (
	"github.com/Dataman-cloud/omega-es/src/config"
	. "github.com/Dataman-cloud/omega-es/src/es"
	_ "github.com/Dataman-cloud/omega-es/src/logger"
	"github.com/Dataman-cloud/omega-es/src/util"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	initGin()
}

func initGin() {
	log.Info("gin starting...")
	addr := config.GetString("host") + ":" + config.GetString("port")
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), optionHandler, authenticate)

	r.GET("/", func(c *gin.Context) {
		c.String(200, "pong")
	})
	es := r.Group("/search")
	{
		es.POST("/index", Search)
		es.POST("/jump", SearchJump)
	}
	s := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func authenticate(ctx *gin.Context) {
	authenticated := false
	author := util.Header(ctx, "Authorization")
	if len(author) > 0 {
		conn := util.Open()
		defer conn.Close()
		uid, err := redis.String(conn.Do("HGET", "s:"+author, "user_id"))
		if err == nil {
			authenticated = true
			ctx.Set("uid", uid)
		} else if err != redis.ErrNil {
			log.Error("[app] got error1 ", err)
		}
	}
	if authenticated {
		ctx.Next()
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": -1})
		ctx.Abort()
	}
}

func optionHandler(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, X-Requested-By, If-Modified-Since, X-File-Name, Cache-Control, X-XSRFToken, Authorization")
	ctx.Header("Content-Type", "application/json")
	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(204)
	}

	ctx.Next()
}
