package main

import (
	_ "github.com/team4yf/fpm-go-plugin-mqtt-client/plugin"
	"github.com/team4yf/fpm-iot-go-middleware/config"
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/pool"
	"github.com/team4yf/fpm-iot-go-middleware/router"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

func init() {

}

var migration model.Migration

func main() {
	app := fpm.New()

	app.AddHook("BEFORE_INIT", func(f *fpm.Fpm) {
		config.Init("")
		// Init the model
		model.CreateDb()
		migration.Install()

		// Init the redis pool
		pool.InitRedis(config.RedisConfig)
	}, 10)

	app.AddHook("AFTER_INIT", func(f *fpm.Fpm) {
		router.LoadPushAPI(app)
		router.LoadDeviceAPI(app)
		router.LoadMQTTUserAPI(app)
	}, 10)

	app.Init()
	// app.Subscribe("$s2d/+/+/send", consumer.DefaultMqttConsumer(app))

	// app.Subscribe("$d2s/+/mcu20/push", consumer.DevicePushConsumer(app))

	app.Run()
}
