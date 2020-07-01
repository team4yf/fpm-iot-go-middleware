package config

import MQTT "github.com/eclipse/paho.mqtt.golang"

//MqttSetting the mqtt config
type MqttSetting struct {
	Options  *MQTT.ClientOptions
	Qos      byte
	Retained bool
}
