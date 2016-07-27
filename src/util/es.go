package util

import (
	"strings"

	"github.com/Dataman-Cloud/omega-es/src/config"
	es "github.com/mattbaird/elastigo/lib"
)

var Conn *es.Conn

var EsSearch = func(index string, _type string, args map[string]interface{}, query interface{}) (es.SearchResult, error) {
	out, err := Conn.Search(index, _type, args, query)
	return out, err
}

func EsInit() {
	Conn = es.NewConn()
	Conn.SetHosts(strings.Split(config.GetConfig().Ec.Hosts, ","))
	Conn.SetPort(config.GetConfig().Ec.Port)
}
