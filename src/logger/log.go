package logger

import (
	"github.com/Dataman-cloud/omega-es/src/config"
	log "github.com/cihub/seelog"
)

var (
	console    bool
	appendfile bool
	file       string
)

func InitLogger() {
	_, console = config.GetStringMapBool("log", "console")
	_, appendfile = config.GetStringMapBool("log", "appendfile")
	_, file = config.GetStringMapString("log", "file")
	if file == "" {
		file = "./log/omega-es.log"
	}
	logger, err := log.LoggerFromConfigAsString(logConfig())
	if err == nil {
		log.ReplaceLogger(logger)
	} else {
		log.Error(err)
	}
}
