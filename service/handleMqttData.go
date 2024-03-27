package service

import (
	"fmt"
	models "goMqtt/database"
)

func VoltageDataProcessing(kind, topic, payload string) {
	db := models.GetDB()
	rows, err := db.Query("SELECT id, mqtt_id, energy FROM chargingpile")
	if err != nil {
		Logger.Error(err)
	} else {
		for rows.Next() {
			var chargingPile models.ChargingPile
			if err := rows.Scan(&chargingPile.Id, &chargingPile.Mqtt_id, &chargingPile.Energy); err != nil {
				Logger.Error(err)
			}
			fmt.Printf("Id:%d, Mqtt_id:%d, Energy:%f\n", chargingPile.Id, chargingPile.Mqtt_id, chargingPile.Energy)
		}
		defer rows.Close()
	}
}

func CurrentDataProcessing(kind, topic, payload string) {
}

func EnergyDataProcessing(kind, topic, payload string) {
}
