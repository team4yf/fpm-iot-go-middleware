package main

import (
	"fmt"
	"strings"

	"github.com/team4yf/fpm-go-pkg/utils"
	"github.com/team4yf/fpm-iot-go-middleware/external/device/env"

	_ "github.com/team4yf/fpm-go-plugin-cache-redis/plugin"
	_ "github.com/team4yf/fpm-go-plugin-mqtt-client/plugin"
	_ "github.com/team4yf/fpm-go-plugin-orm/plugins/pg"
	_ "github.com/team4yf/fpm-go-plugin-tcp/plugin"
	"github.com/team4yf/fpm-iot-go-middleware/consumer"
	"github.com/team4yf/fpm-iot-go-middleware/handler/device"
	"github.com/team4yf/fpm-iot-go-middleware/handler/mqttuser"
	"github.com/team4yf/fpm-iot-go-middleware/handler/project"
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/fpm-iot-go-middleware/router"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

func main() {
	app := fpm.New()

	app.AddHook("BEFORE_INIT", func(f *fpm.Fpm) {
		// Init the model
		dbclient, _ := app.GetDatabase("pg")
		migrator := &model.Migration{
			DS: dbclient,
		}
		migrator.Install()
	}, 10)

	app.AddHook("AFTER_INIT", func(f *fpm.Fpm) {
		router.LoadPushAPI(app)
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

	envDevice := env.NewEnvDevice("jingxun")
	app.Subscribe("#tcp/receive", func(_ string, data interface{}) {
		body := data.(map[string]interface{})
		app.Logger.Debugf("receive tcp data: %X", body["data"])
		buf := body["data"].([]byte)
		envData, err := envDevice.Parse(&buf)
		if err != nil {
			app.Logger.Errorf("parse env tcp data error: %v", err)
			return
		}

		// body["data"] = fmt.Sprintf("%x", body["data"])
		app.Execute("mqttclient.publish", &fpm.BizParam{
			"topic":   fmt.Sprintf(`$d2s/%s/partner/push`, envData.Header.AppID),
			"payload": ([]byte)(utils.JSON2String(envData)),
		})
	})

	mqttuser.InitBiz(app)
	device.InitBiz(app)
	project.InitBiz(app)
	app.Run()
}
