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

var OneMinDataMap = make(map[string]ChargingData)

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
		// // 定義好topic前面四碼去掉後 為chargingPileId
		// chargingPileId := string(needTopicWordRune[4:])

		// if containsChinese(chargingPileId) {
		// 	errMsg := fmt.Sprintf("chargingPileId 解析後含有中文不處理, chargingPileId:%s", chargingPileId)
		// 	return "", errors.New(errMsg)
		// } else {
		// 	return chargingPileId, nil
		// }

		// 205定義好topic去掉最後2碼 為chargingPileId
		chargingPileId := string(needTopicWordRune[:len(needTopicWordRune)-2])
		return chargingPileId, nil
	}

	errMsg := fmt.Sprintf("invalid topic:%s", topic)
	return "", errors.New(errMsg)
}

func handleData(kind, topic, payload string) {
	// log.Printf("處理'%s'數據:%s, 時間:%s, topic:%s\n", kind, payload, time.Now().Format("2006-01-02 15:04:05"), topic)
	Logger.Info("處理", kind, "數據:", payload, ", topic:", topic)

	chargingPileId, err := getChargingPileId(topic)
	if err != nil {
		Logger.Error(err)
	} else {
		// 如果 OneMinDataMap[chargingPileId] 還未初始化，則初始化為空的 ChargingData struct
		if _, ok := OneMinDataMap[chargingPileId]; !ok {
			OneMinDataMap[chargingPileId] = ChargingData{}
		}

		// 取出 chargingData
		chargingData := OneMinDataMap[chargingPileId]

		// 根據 kind 將 payload 添加到ㄇ對應的欄位中
		switch kind {
		case "Voltage":
			chargingData.Voltage = append(chargingData.Voltage, payload)
			// VoltageDataProcessing("Voltage", topic, payload)
		case "Current":
			chargingData.Current = append(chargingData.Current, payload)
			// CurrentDataProcessing("Current", topic, payload)
		case "Energy":
			chargingData.Energy = append(chargingData.Energy, payload)
			// EnergyDataProcessing("Energy", topic, payload)
		default:
			Logger.Info("未定義種類不處理, kind:", kind)
		}

		// 將修改後的 chargingData 放回 map 中
		OneMinDataMap[chargingPileId] = chargingData
	}
}

func onMessageReceived(message MQTT.Message, input string) {
	topic := message.Topic()
	payload := string(message.Payload())

	switch {
	case input == "電壓" && strings.Contains(topic, "電壓"):
		handleData("Voltage", topic, payload)
	case input == "電流" && strings.Contains(topic, "電流"):
		handleData("Current", topic, payload)
	case input == "用電" && (strings.Contains(topic, "用電") || strings.Contains(topic, "總電")):
		handleData("Energy ", topic, payload)
	case input == "全部":
		if strings.Contains(topic, "電壓") {
			handleData("Voltage", topic, payload)
		} else if strings.Contains(topic, "電流") {
			handleData("Current", topic, payload)
		} else if strings.Contains(topic, "用電") || strings.Contains(topic, "總電") {
			handleData("Energy", topic, payload)
		}
	default:
		Logger.Info("輸入參數:", input, "與topic不合不處理, topic:", topic)
	}
}

func printOneMinDataMap() {
	jsonData, err := json.MarshalIndent(OneMinDataMap, "", "  ")
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return
	}
	Logger.Info("JSON data:", string(jsonData))
	Logger.Info("======================================")
}

func GetRawMqttMain(input string) {
	userName := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	clientId := os.Getenv("CLIENT_ID")
	topic := "EZ01/device/#"
	broker := "wss://ehome.ezcon.com.tw:443/ws"
	// broker := "ws://192.168.0.205:11883/ws"
	// broker := "ws://eztw.in:6190/ws"

	client, err := NewClient(broker, userName, password, clientId, func(message MQTT.Message) {
		// 在匿名函数内部创建闭包，调用 onMessageReceived 函数并传递额外的参数 input
		onMessageReceived(message, input)
	})
	if err != nil {
		Logger.Error("Error get newClient:", err)
		return
	}

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		Logger.Error("Error subscribing to MQTT topic:", token.Error())
		return
	}
	Logger.Info("已連接到 MQTT 代理,並訂閱了主題:", topic)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// 建立定時器，每一分鐘觸發一次印出 OneMinDataMap 的操作
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

	defer client.Disconnect(250)
}
