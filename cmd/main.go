package main

import (
	"fmt"

	config "github.com/team4yf/fpm-iot-go-middleware/config"
	"github.com/team4yf/fpm-iot-go-middleware/internal/core"
	s "github.com/team4yf/fpm-iot-go-middleware/internal/service"
	"github.com/team4yf/fpm-iot-go-middleware/pkg"
	"github.com/team4yf/fpm-iot-go-middleware/router"
)

func init() {

}

func main() {
	cfg := &config.Config{}

	pubSub := GetPubSub()
	service := GetService()
	app := &core.App{}
	app.Config = cfg
	app.Init(pubSub, service)
	router.Load(app)
	app.Run(fmt.Sprintf("%v:%v",
		cfg.GetConfigOrDefault("server.host", "0.0.0.0"), cfg.GetConfigOrDefault("server.port", "9000")))
}

func GetPubSub() pkg.PubSub {
	return pkg.NewMQTTPubSub(config.MQTT_URL, config.MQTT_USER, config.MQTT_PASS, (byte)(config.MQTT_QOS), config.MQTT_RETAIN)
}

func GetService() s.DeviceService {
	return s.NewSimpleDeviceService(fmt.Sprintf("%s:%s", config.REDIS_HOST, config.REDIS_PORT), config.REDIS_PASS, config.REDIS_DB)
}
