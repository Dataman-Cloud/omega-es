package logger

import (
//log "github.com/cihub/seelog"
)

func logConfig() string {
	logconfig := `<seelog type="asynctimer" asyncinterval="5000000" minlevel="debug">
	            <outputs formatid="main">`
	if console {
		logconfig += `<console/>`
	}
	if appendfile {
		logconfig += `<buffered size="10000" flushperiod="1000">`
		logconfig += `<rollingfile type="size" filename="` + file + `" maxsize="5000000" maxrolls="30" />`
		lgoconfig += `</buffered>`
	}
	logconfig += `</outputs>
	           <formats>
	                <format id="main" format="%Date(2006-01-02 15:04:05Z07:00) [%LEVEL] %Msg%n"/>
	           </formats>
               </seelog>`
	return logconfig
}
