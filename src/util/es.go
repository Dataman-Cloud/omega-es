package util

import (
	"strings"

	"github.com/Dataman-Cloud/omega-es/src/config"
	es "github.com/mattbaird/elastigo/lib"
)

var Conn *es.Conn

func EsInit() {
	Conn = es.NewConn()
	Conn.SetHosts(strings.Split(config.GetConfig().Ec.Hosts, ","))
	Conn.SetPort(config.GetConfig().Ec.Port)
}
