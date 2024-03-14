package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/joho/godotenv/autoload"
)

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	topic := message.Topic()
	payload := string(message.Payload())
	fmt.Printf("收到來自主題 %s 的消息：%s\n", topic, payload)
}

func main() {
	userName := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	clientId := os.Getenv("CLIENT_ID")

	opts := MQTT.NewClientOptions().
		AddBroker("ws://eztw.in:6190/ws").
		SetUsername(userName).
		SetPassword(password).
		SetClientID(clientId).
		SetKeepAlive(2 * time.Second).
		SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
			onMessageReceived(client, msg)
		})

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)

	topic := "#" // 訂閱所有主題
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Printf("已連接到 MQTT 代理，並訂閱了主題 %s\n", topic)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	fmt.Println("接收到結束信號，程序退出")
}
