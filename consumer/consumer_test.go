package consumer

import (
	"testing"
)

func TestSubstribe(t *testing.T) {

	// app.Subscribe("$s2d/+/+/send", DefaultMqttConsumer)
	DefaultMqttConsumer("$s2d/+/+/send", `{"header":{"v":10,"ns":"FPM.Lamp.Light","name":"GroupControl","projId":1,"appId":"ceaa191a","source":"MQTT"},"bind":{},"payload":[{"msgId":"1","cgi":"1","netId":"1","device":{"id":"866971039105809","type":"light","name":"deng1","brand":"lt10","v":"v10","x":{}},"cmd":"TurnOff","arg":{},"feedback":0,"timestamp":1594829482984}]}`)
	// app.Run(fmt.Sprintf("%v:%v",
	// 	config.GetConfigOrDefault("server.host", "0.0.0.0"), "7070"))

}
