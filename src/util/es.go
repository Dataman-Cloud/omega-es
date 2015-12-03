package util

import (
	"github.com/Dataman-cloud/omega-es/src/config"
	log "github.com/cihub/seelog"
	es "github.com/mattbaird/elastigo/lib"
	"strings"
)

var Conn *es.Conn

func InitES() {
	err, hosts := config.GetStringMapString("es", "hosts")
	if err != nil {
		hosts = "localhost"
		log.Warn(err)
	}
	err, port := config.GetStringMapString("es", "port")
	if err != nil {
		port = "9200"
		log.Warn(err)
	}
	Conn = es.NewConn()
	Conn.SetHosts(strings.Split(hosts, ","))
	Conn.SetPort(port)
}
