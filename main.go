package main

import (
	"fmt"
	models "goMqtt/database"
	routes "goMqtt/routes"
	service "goMqtt/service"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// 初始化資料庫
	models.InitDB()
	defer models.CloseDB()

	// 獲取用戶輸入
	input := service.GetUserInput()

	// 初始化日志
	service.InitLogger()

	// 設置路由
	// gin.SetMode(gin.ReleaseMode)
	router := routes.SetupRouter()

	// 啟動http server
	go func() {
		if err := router.Run(":8080"); err != nil {
			fmt.Println("HTTP server 啟動失敗:", err)
		}
	}()

	// 啟動mqtt
	go service.GetRawMqttMain(input)

	// 等待接收到中斷訊號
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	// 接收到结束信號關閉HTTP
	fmt.Println("接收到结束信號,關閉HTTP server")
}
