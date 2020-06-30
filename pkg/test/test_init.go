//this just for test.go init
package test

import (
	"github.com/team4yf/fpm-iot-go-middleware/config"
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/pool"
)

func InitTestConfig(testFile string) {

	config.Init(testFile)
	model.CreateDb()
	migration := &model.Migration{}
	migration.Install()
	pool.Init(config.RedisConfig)
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
