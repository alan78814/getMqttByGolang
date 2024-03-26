package main

import (
	models "goMqtt/database"
	service "goMqtt/service"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	models.InitDB()
	defer models.CloseDB()

	service.InitLogger()
	service.GetRawMqttMain()
}
