package main

import (
	"fmt"
	"log"

	config "github.com/team4yf/fpm-iot-go-middleware/config"
	lintaiv10 "github.com/team4yf/fpm-iot-go-middleware/external/device/light/lintai/v10"
	"github.com/team4yf/fpm-iot-go-middleware/internal/core"
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	s "github.com/team4yf/fpm-iot-go-middleware/internal/service"
	"github.com/team4yf/fpm-iot-go-middleware/pkg"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/pool"
	"github.com/team4yf/fpm-iot-go-middleware/router"
)

func init() {

}

var migration model.Migration

func main() {
	config.Init("")
	// Init the model
	model.CreateDb()
	migration.Install()

	// Init the redis pool
	pool.Init(config.RedisConfig)

	options := &lintaiv10.Options{
		AppID:       "LT0314fbf27a4d2986",
		AppSecret:   "1bc7b874c74623298a6",
		Username:    "18796664408",
		TokenExpire: 60 * 1000 * 24 * 7,

		Enviroment: "prod",
		BaseURL:    "http://101.132.142.5:8088/api",
	}
	client := lintaiv10.NewClient(options)

	err := client.Init()
	if err != nil {
		log.Fatal(err)
	}

	pubSub := GetPubSub()
	service := GetService()
	app := &core.App{}

	cfg := &config.Config{}
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
	return s.NewSimpleDeviceService()
}
