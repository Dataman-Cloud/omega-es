package util

import (
	"github.com/Dataman-Cloud/omega-es/src/config"
	es "github.com/mattbaird/elastigo/lib"
	"strings"
)

var Conn *es.Conn

func init() {
	Conn = es.NewConn()
	Conn.SetHosts(strings.Split(config.GetConfig().Ec.Hosts, ","))
	Conn.SetPort(config.GetConfig().Ec.Port)
}
