package config

import (
	log "github.com/cihub/seelog"
	flags "github.com/jessevdk/go-flags"
)

type LogOptions struct {
}

type Options struct {
	Host       string `short:"h" long:"host" description:"listen address" optional:"yes" default:"0.0.0.0"`
	Port       int    `short:"p" long:"port" description:"listen port" optional:"yes" default:"9200"`
	Console    string `short:"c" long:"console" description:"console log" default:"false"`
	Appendfile string `short:"a" long:"appendfile" description:"append file log" default:"false"`
	File       string `short:"f" long:"file" description:"append file log" default:"./log/omega-es.log"`
}

var options Options

func InitFlag() {
	parse := flags.NewParser(&options, flags.Default)
	_, err := parse.Parse()
	if err != nil {
		log.Error("flag parse error:", err)
	}
}
