package logger

import (
	"github.com/Dataman-Cloud/omega-es/src/config"
	log "github.com/cihub/seelog"
)

var (
	console    bool
	appendfile bool
	file       string
)

func init() {
	console = config.GetConfig().Lc.Console
	appendfile = config.GetConfig().Lc.AppendFile
	file = config.GetConfig().Lc.File
	logger, err := log.LoggerFromConfigAsString(logConfig())
	if err == nil {
		log.ReplaceLogger(logger)
	} else {
		log.Error(err)
	}
}
