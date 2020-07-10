package pubsub

// 用于 mqtt pub/sub 的函数
// 默认使用MQTT的实现,后面会根据情况加入 Kafka 和 Rabbit 的实现

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/team4yf/fpm-iot-go-middleware/config"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/log"
)

//PubSub 定义接口
// 主要包含发布和订阅
type PubSub interface {
	Publish(topic string, payload []byte)
	Subscribe(topic string, handler func(topic, payload interface{}))
}

//mqttPS 定义MQTT 的结构体
// 包含一个 MQTT 的客户端和一些配置信息
type mqttPS struct {
	mClient MQTT.Client
	config  *config.MqttSetting
}

//NewMQTTPubSub 构建实例的函数,用于返回一个MQTT的对象,通过 PubSub 接口返回
func NewMQTTPubSub(c *config.MqttSetting) PubSub {
	client := MQTT.NewClient(c.Options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	pubSub := &mqttPS{
		mClient: client,
		config:  c,
	}
	return pubSub
}

//Publish 实现Publish函数
func (m *mqttPS) Publish(topic string, payload []byte) {
	log.Infof("topic: %s, payload: %s", topic, payload)
	token := m.mClient.Publish(topic, m.config.Qos, m.config.Retained, payload)
	token.Wait()
}

//Subscribe 实现Subscribe
func (m *mqttPS) Subscribe(topic string, handler func(topic, payload interface{})) {
	m.mClient.Subscribe(topic, m.config.Qos, func(_ MQTT.Client, message MQTT.Message) {
		handler(topic, message.Payload())
	})
}
