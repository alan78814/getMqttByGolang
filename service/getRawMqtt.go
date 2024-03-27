package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unicode"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/joho/godotenv/autoload"
)

type ChargingData struct {
	Voltage []string
	Current []string
	Energy  []string
	Status  []string
}

var oneMinDataMap = make(map[string]ChargingData)

func containsChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func getChargingPileId(topic string) (string, error) {
	// topic範例"EZ01/device/總表用電101" 用/區分取第三個
	topicSeparate := strings.Split(topic, "/")
	needTopicWordRune := []rune(topicSeparate[2])

	if len(topicSeparate) >= 3 {
		// 定義好topic前面四碼去掉後 為chargingPileId
		chargingPileId := string(needTopicWordRune[4:])

		if containsChinese(chargingPileId) {
			errMsg := fmt.Sprintf("chargingPileId 解析後含有中文不處理, chargingPileId:%s", chargingPileId)
			return "", errors.New(errMsg)
		} else {
			return chargingPileId, nil
		}
	}

	errMsg := fmt.Sprintf("invalid topic:%s", topic)
	return "", errors.New(errMsg)
}

func handleData(kind, topic, payload string) {
	// log.Printf("處理'%s'數據:%s, 時間:%s, topic:%s\n", kind, payload, time.Now().Format("2006-01-02 15:04:05"), topic)
	Logger.Info("處理", kind, "數據:", payload, " topic:", topic)

	chargingPileId, err := getChargingPileId(topic)
	if err != nil {
		Logger.Error(err)
	} else {
		// 如果 oneMinDataMap[chargingPileId] 還未初始化，則初始化為空的 ChargingData struct
		if _, ok := oneMinDataMap[chargingPileId]; !ok {
			oneMinDataMap[chargingPileId] = ChargingData{}
		}

		// 取出 chargingData
		chargingData := oneMinDataMap[chargingPileId]

		// 根據 kind 將 payload 添加到ㄇ對應的欄位中
		switch kind {
		case "Voltage":
			chargingData.Voltage = append(chargingData.Voltage, payload)
			VoltageDataProcessing("Voltage", topic, payload)
		case "Current":
			chargingData.Current = append(chargingData.Current, payload)
			CurrentDataProcessing("Current", topic, payload)
		case "Energy":
			chargingData.Energy = append(chargingData.Energy, payload)
			EnergyDataProcessing("Energy", topic, payload)
		default:
			Logger.Info("未定義種類不處理, kind:", kind)
		}

		// 將修改後的 chargingData 放回 map 中
		oneMinDataMap[chargingPileId] = chargingData
	}
}

func onMessageReceived(message MQTT.Message) {
	topic := message.Topic()
	payload := string(message.Payload())

	switch {
	case strings.Contains(topic, "電壓"):
		handleData("Voltage", topic, payload)
	case strings.Contains(topic, "電流"):
		handleData("Current", topic, payload)
	case strings.Contains(topic, "用電"):
		handleData("Energy ", topic, payload)
	default:
		Logger.Info("未定義主題不處理, topic:", topic)
	}
}

func printOneMinDataMap() {
	jsonData, err := json.MarshalIndent(oneMinDataMap, "", "  ")
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return
	}
	Logger.Info("JSON data:", string(jsonData))
	Logger.Info("======================================")
}

func GetRawMqttMain() {
	userName := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	clientId := os.Getenv("CLIENT_ID")
	topic := "EZ01/device/#"
	broker := "ws://192.168.0.208:11883/ws"
	// broker := "ws://eztw.in:6190/ws"

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
		Logger.Error("Error connecting to MQTT broker:", token.Error())
		return // 終止函式執行
	}
	defer client.Disconnect(250)

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		Logger.Error("Error subscribing to MQTT topic:", token.Error())
		return // 終止函式執行
	}
	Logger.Info("已連接到 MQTT 代理,並訂閱了主題:", topic)

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
