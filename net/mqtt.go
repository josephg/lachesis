package net

import (
	"encoding/json"
	"github.com/andrecronje/lachesis/utils"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MqttSocket mqttt socket connection for communication
type MqttSocket struct {
	options *MQTT.ClientOptions
	client  MQTT.Client
}

// NewMqttSocket returns new MqttSocket
func NewMqttSocket(host string, callback MQTT.MessageHandler) *MqttSocket {
	options := MQTT.NewClientOptions().AddBroker(host)
	options.AutoReconnect = true
	options.OnConnect = func(client MQTT.Client) {
		// MQTT client connected to server
	}
	options.OnConnectionLost = func(client MQTT.Client, e error) {
		// MQTT client connection lost with server
	}
	options.SetClientID(utils.NewUUID())
	options.SetDefaultPublishHandler(callback)
	return &MqttSocket{
		options: options,
	}
}

// Connect creates connection to server or returns error if fails
func (ms *MqttSocket) Connect() error {
	ms.client = MQTT.NewClient(ms.options)
	if token := ms.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// FireEvent publish event to a specific topic or returns error if fails
func (ms *MqttSocket) FireEvent(v interface{}, topic string) error {
	v, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if token := ms.client.Publish(topic, 2, false, v); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Listen subscribes to a specific topic to get published event on the topic
func (ms *MqttSocket) Listen(topic string) error {
	if token := ms.client.Subscribe(topic, 2, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Disconnect disconnect client from server
func (ms *MqttSocket) Disconnect() {
	ms.client.Disconnect(0)
}
