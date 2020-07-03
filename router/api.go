package router

import (
	"github.com/team4yf/fpm-iot-go-middleware/handler/device"
	"github.com/team4yf/fpm-iot-go-middleware/internal/core"
)

func Load(app *core.App) {
	app.Post("/push/{device}/{brand}/{event}", device.PushHandler(app))
}
