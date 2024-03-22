package main

import (
	"goMqtt/service"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	service.InitLogger()
	service.GetRawMqttMain()
}
