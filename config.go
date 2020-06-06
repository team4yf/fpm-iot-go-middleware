package main

import "fmt"

var (
	PORT              = "3009"
	REDIS_HOST        = "localhost"
	REDIS_PORT        = "6379"
	REDIS_DB          = 13
	REDIS_PASS        = "admin123"
	REDIS_PREFIX      = "drm"
	MQTT_URL          = "www.ruichen.top:1883"
	MQTT_USER         = "admin"
	MQTT_PASS         = "123123123"
	MQTT_EVENT_PREFIX = "^drm"
	MQTT_EVENT_QOS    = 0
	MQTT_EVENT_RETAIN = false
)

type Config struct{}

func (c *Config) GetPubSub() PubSub {
	return NewMQTTPubSub(MQTT_URL, MQTT_USER, MQTT_PASS, (byte)(MQTT_EVENT_QOS), MQTT_EVENT_RETAIN)
}

func (c *Config) GetService() Service {
	return NewRedisService(fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT), REDIS_PASS, REDIS_DB)
}
