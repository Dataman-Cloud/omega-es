package config

import (
	"errors"
	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("omega-es")
	viper.AddConfigPath("./")
	viper.AddConfigPath("$HOME/.omega/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("can't read config file:", err)
	}
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
