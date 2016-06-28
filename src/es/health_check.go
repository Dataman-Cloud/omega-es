package es

import (
	"fmt"
	"github.com/Dataman-Cloud/health_checker"
	"github.com/Dataman-Cloud/omega-es/src/config"
	"github.com/Dataman-Cloud/omega-es/src/util"
	"github.com/gin-gonic/gin"
)

func Health(ctx *gin.Context) {
	checker := health_checker.NewHealthChecker("omega-billing")
	conf := config.GetConfig()
	redisDsn := fmt.Sprintf("%s:%d",
		conf.Rc.Host, conf.Rc.Port)
	checker.AddCheckPoint("redis", redisDsn, nil, nil)

	mysqlDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		conf.Mc.UserName, conf.Mc.PassWord, conf.Mc.Host, conf.Mc.Port, conf.Mc.DataBase)
	checker.AddCheckPoint("mysql", mysqlDsn, nil, nil)
	/*
		mqDsn := fmt.Sprintf("amqp://%s:%s@%s:%d/",
			conf.Mq.User, conf.Mq.PassWord, conf.Mq.Host, conf.Mq.Port)
		checker.AddCheckPoint("mq", mqDsn, nil, nil)
	*/
	util.ReturnOKGin(ctx, checker.Check())
}
