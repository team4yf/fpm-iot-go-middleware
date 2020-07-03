package main

import (
	"fmt"

	config "github.com/team4yf/fpm-iot-go-middleware/config"
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
	//TODO: get all rest client config
	router.Load(app)
	app.Run(fmt.Sprintf("%v:%v",
		config.GetConfigOrDefault("server.host", "0.0.0.0"), config.GetConfigOrDefault("server.port", "9000")))
}