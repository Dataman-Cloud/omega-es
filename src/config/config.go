package config

import (
	"bufio"
	"errors"
	log "github.com/cihub/seelog"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func init() {
	InitConfig("deploy/env")
}

type Config struct {
	Host    string
	Port    uint16
	Murl    string
	Userurl string
	Appurl  string
	Lc      LogConfig       `mapstructure:"log"`
	Ec      EsConfig        `mapstructure:"es"`
	Rc      RedisConfig     `mapstructure:"redis"`
	Mc      MysqlConfig     `mapstructure:"mysql"`
	Ch      ChronosConfig   `mapstructure:"chronos"`
	Sh      SchedulerConfig `mapstructure:"scheduler"`
}

type LogConfig struct {
	Console    bool
	AppendFile bool
	File       string
	Level      string
	Formatter  string
	MaxSize    uint64
}

type EsConfig struct {
	Hosts string
	Port  string
}

type RedisConfig struct {
	Host string
	Port uint16
}

type MysqlConfig struct {
	Host         string
	Port         uint16
	MaxIdleConns int64
	MaxOpenConns int64
	DataBase     string
	UserName     string
	PassWord     string
}

type ChronosConfig struct {
	Host string
	Port uint16
}

type SchedulerConfig struct {
	Host string
	Port uint16
}

type EnvEntry struct {
	OMEGAES_NET_HOST           string `required:"true"`
	OMEGAES_NET_PORT           uint16 `required:"true"`
	OMEGAES_NET_APPURL         string `required:"true"`
	OMEGAES_LOG_CONSOLE        bool   `required:"true"`
	OMEGAES_LOG_APPENDFILE     bool   `required:"true"`
	OMEGAES_LOG_FILE           string `required:"false"`
	OMEGAES_LOG_LEVEL          string `required:"true"`
	OMEGAES_LOG_FORMATTER      string `required:"true"`
	OMEGAES_LOG_MAXSIZE        uint64 `required:"true"`
	OMEGAES_ES_HOSTS           string `required:"true"`
	OMEGAES_ES_PORT            string `required:"true"`
	OMEGAES_REDIS_HOST         string `required:"true"`
	OMEGAES_REDIS_PORT         uint16 `required:"true"`
	OMEGAES_MYSQL_HOST         string `required:"true"`
	OMEGAES_MYSQL_PORT         uint16 `required:"true"`
	OMEGAES_MYSQL_MAXIDLECONNS int64  `required:"true"`
	OMEGAES_MYSQL_MAXOPENCONNS int64  `required:"true"`
	OMEGAES_MYSQL_DATABASE     string `required:"true"`
	OMEGAES_MYSQL_USERNAME     string `required:"true"`
	OMEGAES_MYSQL_PASSOWRD     string `required:"true"`
}

var config Config

func GetConfig() Config {
	return config
}

func InitConfig(path string) {
	loadEnvFile(path)

	envEntry := NewEnvEntry()
	config.Host = envEntry.OMEGAES_NET_HOST
	config.Port = envEntry.OMEGAES_NET_PORT
	config.Appurl = envEntry.OMEGAES_NET_APPURL
	config.Lc.Console = envEntry.OMEGAES_LOG_CONSOLE
	config.Lc.AppendFile = envEntry.OMEGAES_LOG_APPENDFILE
	config.Lc.File = envEntry.OMEGAES_LOG_FILE
	config.Lc.Level = envEntry.OMEGAES_LOG_LEVEL
	config.Lc.Formatter = envEntry.OMEGAES_LOG_FORMATTER
	config.Lc.MaxSize = envEntry.OMEGAES_LOG_MAXSIZE
	config.Ec.Hosts = envEntry.OMEGAES_ES_HOSTS
	config.Ec.Port = envEntry.OMEGAES_ES_PORT
	config.Rc.Host = envEntry.OMEGAES_REDIS_HOST
	config.Rc.Port = envEntry.OMEGAES_REDIS_PORT
	config.Mc.Host = envEntry.OMEGAES_MYSQL_HOST
	config.Mc.Port = envEntry.OMEGAES_MYSQL_PORT
	config.Mc.MaxIdleConns = envEntry.OMEGAES_MYSQL_MAXIDLECONNS
	config.Mc.MaxOpenConns = envEntry.OMEGAES_MYSQL_MAXOPENCONNS
	config.Mc.DataBase = envEntry.OMEGAES_MYSQL_DATABASE
	config.Mc.UserName = envEntry.OMEGAES_MYSQL_USERNAME
	config.Mc.PassWord = envEntry.OMEGAES_MYSQL_PASSOWRD
}

func NewEnvEntry() *EnvEntry {
	envEntry := &EnvEntry{}

	val := reflect.ValueOf(envEntry).Elem()

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		required := typeField.Tag.Get("required")

		env := os.Getenv(typeField.Name)

		if env == "" && required == "true" {
			exitMissingEnv(typeField.Name)
		}

		var envEntryValue interface{}
		var err error
		valueFiled := val.Field(i).Interface()
		value := val.Field(i)
		switch valueFiled.(type) {
		case int64:
			envEntryValue, err = strconv.ParseInt(env, 10, 64)

		case int16:
			envEntryValue, err = strconv.ParseInt(env, 10, 16)
			_, ok := envEntryValue.(int64)
			if !ok {
				exitCheckEnv(typeField.Name, err)
			}
			envEntryValue = int16(envEntryValue.(int64))
		case uint16:
			envEntryValue, err = strconv.ParseUint(env, 10, 16)

			_, ok := envEntryValue.(uint64)
			if !ok {
				exitCheckEnv(typeField.Name, err)
			}
			envEntryValue = uint16(envEntryValue.(uint64))
		case uint64:
			envEntryValue, err = strconv.ParseUint(env, 10, 64)
		case bool:
			envEntryValue, err = strconv.ParseBool(env)
		default:
			envEntryValue = env
		}

		if err != nil {
			exitCheckEnv(typeField.Name, err)
		}
		value.Set(reflect.ValueOf(envEntryValue))
	}

	return envEntry
}

func loadEnvFile(envfile string) {
	// load the environment file
	f, err := os.Open(envfile)
	if err == nil {
		defer f.Close()

		r := bufio.NewReader(f)
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				break
			}

			key, val, err := parseln(string(line))
			if err != nil {
				continue
			}

			if len(os.Getenv(strings.ToUpper(key))) == 0 {
				err1 := os.Setenv(strings.ToUpper(key), val)
				if err1 != nil {
					log.Error(err1.Error())
				}
			}
		}
	}
}

// helper function to parse a "key=value" environment variable string.
func parseln(line string) (key string, val string, err error) {
	line = removeComments(line)
	if len(line) == 0 {
		return
	}
	splits := strings.SplitN(line, "=", 2)

	if len(splits) < 2 {
		err = errors.New("missing delimiter '='")
		return
	}

	key = strings.Trim(splits[0], " ")
	val = strings.Trim(splits[1], ` "'`)
	return

}

// helper function to trim comments and whitespace from a string.
func removeComments(s string) (_ string) {
	if len(s) == 0 || string(s[0]) == "#" {
		return
	} else {
		index := strings.Index(s, " #")
		if index > -1 {
			s = strings.TrimSpace(s[0:index])
		}
	}
	return s
}

func exitMissingEnv(env string) {
	log.Errorf("program exit missing config for env %s", env)
	log.Flush()
	os.Exit(1)
}

func exitCheckEnv(env string, err error) {
	log.Errorf("Check env %s, %s", env, err.Error())
}
