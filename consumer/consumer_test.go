package consumer

import (
	"testing"
)

func TestSubstribe(t *testing.T) {

	// config.Init("../conf/config.test.json")
	// // Init the model
	// model.CreateDb()

	// // Init the redis pool
	// pool.InitRedis(config.RedisConfig)

	// app := &core.App{}

	// app.Init()

	// handler := DefaultMqttConsumer(app)
	// app.Subscribe("$s2d/+/+/send", handler)
	// // app.Subscribe("$s2d/+/+/send", DefaultMqttConsumer)
	// handler("$s2d/+/+/send", ([]byte)(`{"header":{"v":10,"ns":"FPM.Lamp.Light","name":"Control","projId":1,"appId":"ceaa191a","source":"MQTT"},"bind":{},"payload":[{"msgId":"1","cgi":"1","netId":"1","device":{"id":"866971039105809","type":"light","name":"deng1","brand":"lt10","v":"v10","x":{}},"cmd":"command","arg":[{"circuit":1,"imei":"866971039105809","commandType":"BRIGHTNESS","commandValue":30,"type":2}],"feedback":1,"timestamp":1594829482984}]}`))
	// app.Run(fmt.Sprintf("%v:%v",
	// 	config.GetConfigOrDefault("server.host", "0.0.0.0"), "7070"))

}
