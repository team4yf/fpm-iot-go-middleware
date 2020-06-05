package main

// 用于 mqtt pub/sub 的函数
import (
	"fmt"
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type PubSub interface {
	Publish(topic string, payload interface{})
	Subscribe(topic string, handler func(topic, payload interface{}))
}

type MQTTPubSub struct {
	mClient  MQTT.Client
	qos      byte
	retained bool
}

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

func (m *MQTTPubSub) Publish(topic string, payload interface{}) {
	log.Printf("topic: %s, payload: %s", topic, payload)
	token := m.mClient.Publish(topic, m.qos, m.retained, payload)
	token.Wait()
}

func (m *MQTTPubSub) Subscribe(topic string, handler func(topic, payload interface{})) {

}
