//Package router the api router
package router

import (
	"github.com/team4yf/fpm-iot-go-middleware/handler/device"
	"github.com/team4yf/fpm-iot-go-middleware/internal/core"
)

//LoadPushAPI 第三方设备平台推送过来的设备信息
func LoadPushAPI(app *core.App) {
	app.Post("/push/{device}/{brand}/{event}", device.PushHandler(app))

}

//LoadDeviceAPI 设备管理相关的接口
func LoadDeviceAPI(app *core.App) {
	app.Post("/device/{type}/{brand}/create", device.CreateHandler(app))
}
