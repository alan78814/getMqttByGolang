package service

import (
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type CallbackFunc func(msg MQTT.Message)

func NewClient(broker, userName, password, clientId string, callbackFunc CallbackFunc) (MQTT.Client, error) {
	opts := MQTT.NewClientOptions().
		AddBroker(broker).
		SetUsername(userName).
		SetPassword(password).
		SetClientID(clientId).
		SetKeepAlive(2 * time.Second).
		SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
			callbackFunc(msg)
		})
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}
