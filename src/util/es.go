package util

import (
	"github.com/Dataman-Cloud/omega-es/src/config"
	log "github.com/cihub/seelog"
	es "github.com/mattbaird/elastigo/lib"
	"strings"
)

var Conn *es.Conn

func init() {
	err, hosts := config.GetStringMapString("es", "hosts")
	if err != nil {
		log.Error(err)
	}
	err, port := config.GetStringMapString("es", "port")
	if err != nil {
		port = "9200"
		log.Warn("can't find es port default:9200")
	}
	Conn = es.NewConn()
	Conn.SetHosts(strings.Split(hosts, ","))
	Conn.SetPort(port)
}
