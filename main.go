package main

import (
	"fmt"

	"github.com/team4yf/fpm-iot-go-middleware/config"
	"github.com/team4yf/fpm-iot-go-middleware/consumer"
	"github.com/team4yf/fpm-iot-go-middleware/internal/core"
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
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
	pool.InitRedis(config.RedisConfig)

	app := &core.App{}

	app.Init()

	router.LoadPushAPI(app)
	router.LoadDeviceAPI(app)

	app.Subscribe("$s2d/+/+/send", consumer.DefaultMqttConsumer(app))

	app.Run(fmt.Sprintf("%v:%v",
		config.GetConfigOrDefault("server.host", "0.0.0.0"), config.GetConfigOrDefault("server.port", "9000")))
}
