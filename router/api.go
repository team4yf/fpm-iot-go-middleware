//Package router the api router
package router

import (
	"github.com/team4yf/fpm-iot-go-middleware/handler/device"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

//LoadPushAPI 第三方设备平台推送过来的设备信息
func LoadPushAPI(fpmApp *fpm.Fpm) {
	device.Init()
	fpmApp.BindHandler("/push/{device}/{brand}/{event}", device.PushHandler()).Methods("POST")
}
