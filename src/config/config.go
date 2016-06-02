package config

import (
	"errors"
	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

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
	MaxSize    uint32
}

type EsConfig struct {
	Hosts string
	Port  uint16
}

type RedisConfig struct {
	Host string
	Port uint16
}

type MysqlConfig struct {
	Host         string
	Port         uint16
	MaxIdleConns uint8
	MaxOpenConns uint8
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

var pairs Config

func init() {
	viper.SetConfigName("omega-es")
	viper.AddConfigPath("./")
	viper.AddConfigPath("/etc/omega/")
	viper.AddConfigPath("$HOME/.omega/")
	viper.AddConfigPath("/")
	if err := viper.ReadInConfig(); err != nil {
		log.Error("can't read config file:", err)
	}

	if err := viper.Unmarshal(&pairs); err != nil {
		log.Errorf("unmarshal config to struct error: %v", err)
	}
}

func GetConfig() Config {
	return pairs
}

func Get(key string) interface{} {
	return viper.Get(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func GetStringMapBool(key, mkey string) (error, bool) {
	m := viper.GetStringMap(key)[mkey]
	if m != nil {
		return nil, m.(bool)
	}
	return errors.New("can't found:" + key + "," + mkey), false
}

func GetStringMapString(key, mkey string) (error, string) {
	m := viper.GetStringMapString(key)[mkey]
	if m != "" {
		return nil, m
	}
	return errors.New("can't found:" + key + "," + mkey), ""
}

func GetStringMapInt(key, mkey string) (error, int) {
	m := viper.GetStringMap(key)[mkey]
	if m != "" {
		return nil, m.(int)
	}
	return errors.New("can't found:" + key + "," + mkey), 0
}
