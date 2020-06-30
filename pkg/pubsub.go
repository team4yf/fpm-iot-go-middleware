package pkg

// 用于 mqtt pub/sub 的函数
// 默认使用MQTT的实现,后面会根据情况加入 Kafka 和 Rabbit 的实现

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/log"
)

// 定义接口
// 主要包含发布和订阅
type PubSub interface {
	Publish(topic string, payload interface{})
	Subscribe(topic string, handler func(topic, payload interface{}))
}

// 定义MQTT 的结构体
// 包含一个 MQTT 的客户端和一些配置信息
type MQTTPubSub struct {
	mClient  MQTT.Client
	qos      byte
	retained bool
}

// 构建实例的函数,用于返回一个MQTT的对象,通过 PubSub 接口返回
func NewMQTTPubSub(url, user, pass string, qos byte, retained bool) PubSub {
	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", url))
	// opts.SetClientID("go-simple")
	opts.SetUsername(user)
	opts.SetPassword(pass)
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error().Error())
	}
	pubSub := &MQTTPubSub{
		mClient:  client,
		qos:      qos,
		retained: retained,
	}
	return pubSub
}

// 实现Publish函数
func (m *MQTTPubSub) Publish(topic string, payload interface{}) {
	log.Infof("topic: %s, payload: %s", topic, payload)
	token := m.mClient.Publish(topic, m.qos, m.retained, payload)
	token.Wait()
}

// TODO: 实现Subscribe
func (m *MQTTPubSub) Subscribe(topic string, handler func(topic, payload interface{})) {

}
