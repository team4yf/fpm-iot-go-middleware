package main

import (
	"fmt"
	"strings"

	_ "github.com/team4yf/fpm-go-plugin-cache-redis/plugin"
	_ "github.com/team4yf/fpm-go-plugin-mqtt-client/plugin"
	_ "github.com/team4yf/fpm-go-plugin-tcp/plugin"
	"github.com/team4yf/fpm-iot-go-middleware/config"
	"github.com/team4yf/fpm-iot-go-middleware/consumer"
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
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
	}, 10)

	app.AddHook("AFTER_INIT", func(f *fpm.Fpm) {
		router.LoadPushAPI(app)
		router.LoadDeviceAPI(app)
		router.LoadMQTTUserAPI(app)
	}, 10)

	app.Init()

	//执行订阅的函数
	app.Execute("mqttclient.subscribe", &fpm.BizParam{
		"topics": []string{"$s2d/+/+/send", "$d2s/+/mcu20/push"},
	})

	mqttHandler := consumer.DefaultMqttConsumer(app)
	mcuHandler := consumer.DevicePushConsumer(app)
	app.Subscribe("#mqtt/receive", func(_ string, data interface{}) {
		//data 通常是 byte[] 类型，可以转成 string 或者 map
		body := data.(map[string]interface{})
		topic := body["topic"].(string)
		switch {
		case strings.HasSuffix(topic, "send"):
			mqttHandler(topic, body["payload"])
		case strings.HasSuffix(topic, "mcu20/push"):
			mcuHandler(topic, body["payload"])
		}
	})

	app.Subscribe("#tcp/receive", func(_ string, data interface{}) {
		body := data.(map[string]interface{})
		app.Logger.Debugf("data: %X", body["data"])
		body["data"] = fmt.Sprintf("%x", body["data"])
		app.Execute("mqttclient.publish", &fpm.BizParam{
			"topic":   `$d2s/tcp`,
			"payload": ([]byte)(utils.JSON2String(body)),
		})
	})

	app.Run()
}
