package main

import (
	service "goMqtt/service"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	service.GetRawMqttMain()
}
