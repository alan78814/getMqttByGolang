package main

import (
	models "goMqtt/database"
	service "goMqtt/service"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	models.InitDB()
	defer models.CloseDB()

	input := service.GetUserInput()
	service.InitLogger()
	service.GetRawMqttMain(input)
}
