package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/joho/godotenv/autoload"
)

var oneMinDataMap = make(map[string][]string)

func onMessageReceived(message MQTT.Message) {
	topic := message.Topic()
	payload := string(message.Payload())

	if _, ok := oneMinDataMap[topic]; ok {
		// fmt.Printf("主題 %s 存在於 oneMinDataMap 中\n", topic)
		oneMinDataMap[topic] = append(oneMinDataMap[topic], payload)
	} else {
		// fmt.Printf("主題 %s 不存在於 oneMinDataMap 中\n", topic)
		oneMinDataMap[topic] = make([]string, 0)
	}
}

func printOneMinDataMap() {
	fmt.Println("Current oneMinDataMap content:")
	currentTime := time.Now()

	// for key, value := range oneMinDataMap {
	// 	fmt.Printf("Key: %s, Value: %v\n", key, value)
	// }

	jsonData, err := json.MarshalIndent(oneMinDataMap, "", "    ")
	if err != nil {
		fmt.Println("無法將資料轉換為 JSON 格式：", err)
		return
	}

	fmt.Println(string(jsonData))

	fmt.Println("現在的時間是：", currentTime)
	fmt.Printf("======================================")
}

func main() {
	userName := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	clientId := os.Getenv("CLIENT_ID")
	topic := "EZchargingPile/status/EZ208/#"
	broker := "ws://eztw.in:6190/ws"

	opts := MQTT.NewClientOptions().
		AddBroker(broker).
		SetUsername(userName).
		SetPassword(password).
		SetClientID(clientId).
		SetKeepAlive(2 * time.Second).
		SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
			onMessageReceived(msg)
		})

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Printf("已連接到 MQTT 代理，並訂閱了主題 %s\n", topic)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// 建立定時器，每一分鐘觸發一次印出 oneMinDataMap 的操作
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				printOneMinDataMap()
			case <-sig:
				fmt.Println("接收到結束信號，程序退出")
				return
			}
		}
	}()

	<-sig

	fmt.Println("接收到結束信號，程序退出")
}