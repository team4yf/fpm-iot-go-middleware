package config

import (
	"log"

	"github.com/spf13/viper"
)

var (
	REDIS_HOST  string
	REDIS_PORT  string
	REDIS_DB    int
	REDIS_PASS  string
	MQTT_URL    string
	MQTT_USER   string
	MQTT_PASS   string
	MQTT_QOS    int
	MQTT_RETAIN bool
	PG_HOST     string
	PG_PORT     int
	PG_USER     string
	PG_DB       string
	PG_PASS     string
	PG_SHOWSQL  bool
)

func init() {
	viper.SetDefault("redis", map[string]interface{}{"host": "localhost", "port": 6379, "db": 0, "password": "admin123"})
	viper.SetDefault("mqtt", map[string]interface{}{"url": "www.ruichen.top:1883", "user": "admin", "qos": 0, "retain": false, "pass": "123123123"})
	viper.SetDefault("postgresql", map[string]interface{}{"host": "localhost", "port": 5432, "db": "fim", "user": "postgres", "password": "Fim741235896", "showsql": true})
	viper.SetConfigType("json")
	viper.AddConfigPath(".") // 设置配置文件和可执行二进制文件在用一个目录
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("no such config file")
		} else {
			// Config file was found but another error was produced
			log.Println("read config error")
		}
		log.Fatal(err) // 读取配置文件失败致命错误
	}
	REDIS_HOST = viper.GetString("redis.host")
	REDIS_PORT = viper.GetString("redis.port")
	REDIS_DB = viper.GetInt("redis.db")
	REDIS_PASS = viper.GetString("redis.password")

	PG_HOST = viper.GetString("postgresql.host")
	PG_PORT = viper.GetInt("postgresql.port")
	PG_USER = viper.GetString("postgresql.user")
	PG_DB = viper.GetString("postgresql.db")
	PG_PASS = viper.GetString("postgresql.password")
	PG_SHOWSQL = viper.GetBool("postgresql.showsql")

	MQTT_URL = viper.GetString("mqtt.url")
	MQTT_USER = viper.GetString("mqtt.user")
	MQTT_PASS = viper.GetString("mqtt.pass")
	MQTT_QOS = viper.GetInt("mqtt.qos")
	MQTT_RETAIN = viper.GetBool("mqtt.retain")
	log.Println(REDIS_HOST, REDIS_PORT, REDIS_DB, REDIS_PASS, MQTT_URL, MQTT_USER, MQTT_PASS, MQTT_QOS, MQTT_RETAIN)
}

type Config struct{}

func (c *Config) GetConfigOrDefault(key string, dft interface{}) interface{} {
	if viper.IsSet(key) {
		return viper.Get(key)
	}
	return dft

}
func (c *Config) IsSet(key string) bool {
	return viper.IsSet(key)
}
func (c *Config) GetMapOrDefault(key string, dft map[string]interface{}) map[string]interface{} {
	if viper.IsSet(key) {
		return viper.GetStringMap(key)
	}
	return dft
}