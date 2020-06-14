package main

import (
	"fmt"

	config "github.com/team4yf/fpm-iot-go-middleware/config"
	"github.com/team4yf/fpm-iot-go-middleware/internal/core"
	s "github.com/team4yf/fpm-iot-go-middleware/internal/service"
	"github.com/team4yf/fpm-iot-go-middleware/pkg"
)

func init() {

}

func main() {
	cfg := &config.Config{}

	pubSub := GetPubSub()
	service, deviceService := GetService()
	app := &core.App{}
	app.Config = cfg
	app.Init(pubSub, service, deviceService)

	app.Run(fmt.Sprintf("%v:%v",
		cfg.GetConfigOrDefault("server.host", "0.0.0.0"), cfg.GetConfigOrDefault("server.port", "9000")))
}

func GetPubSub() pkg.PubSub {
	return pkg.NewMQTTPubSub(config.MQTT_URL, config.MQTT_USER, config.MQTT_PASS, (byte)(config.MQTT_QOS), config.MQTT_RETAIN)
}

func GetService() (s.Service, s.DeviceService) {
	return s.NewRedisService(fmt.Sprintf("%s:%s", config.REDIS_HOST, config.REDIS_PORT), config.REDIS_PASS, config.REDIS_DB),
		s.NewSimpleDeviceService()
}
