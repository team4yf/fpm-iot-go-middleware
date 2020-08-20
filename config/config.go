package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	//Db the db config
	Db *DBSetting
	//RedisConfig the redis connect config
	RedisConfig *redis.Options
	//AppName the app name
	AppName string
)

// Config 读取配置
type Config struct {
	Name string
}

// Init 初始化配置，默认读取config.local.json
func Init(configFilePath string) error {
	c := Config{
		Name: configFilePath,
	}

	// 初始化配置文件
	if err := c.loadConfig(); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) loadConfig() error {

	AppName = viper.GetString("name")
	Db = &DBSetting{
		Engine:   viper.GetString("db.engine"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		Database: viper.GetString("db.database"),
		Charset:  viper.GetString("db.charset"),
		ShowSQL:  viper.GetBool("db.showSql"),
	}
	RedisConfig = &redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.passwd"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool"),
	}

	return nil
}

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
