package config

import (
	"github.com/spf13/viper"
)

//GetConfigOrDefault get config from the config file or env, return default value if not set
func GetConfigOrDefault(key string, dft interface{}) interface{} {
	if viper.IsSet(key) {
		return viper.Get(key)
	}
	return dft
}

//GetMapOrDefault get a map config from the config file or env, return default value if not set.
func GetMapOrDefault(key string, dft map[string]interface{}) map[string]interface{} {
	if viper.IsSet(key) {
		return viper.GetStringMap(key)
	}
	return dft
}
