package service

import (
	models "goMqtt/database"
)

func VoltageDataProcessing(kind, topic, payload string) {
	db := models.GetDB()
	rows, err := db.Query("SELECT id, mqtt_id, energy FROM chargingpile")
	if err != nil {
		Logger.Error(err)
	} else {
		chargingPileId, err := getChargingPileId(topic)
		if err != nil {
			Logger.Error(err)
		} else {
			if _, ok := OneMinDataMap[chargingPileId]; ok {
				chargingData := OneMinDataMap[chargingPileId]

				if lastStatus := chargingData.Status; len(lastStatus) > 0 {
					lastStatus := lastStatus[len(lastStatus)-1]

					switch lastStatus {
					case "standby":
						newStatus := "standby"
						chargingData.Status = append(chargingData.Status, newStatus)
						OneMinDataMap[chargingPileId] = chargingData
					case "charging":
						newStatus := "charging"
						chargingData.Status = append(chargingData.Status, newStatus)
						OneMinDataMap[chargingPileId] = chargingData
					default:
						newStatus := "standby"
						chargingData.Status = append(chargingData.Status, newStatus)
						OneMinDataMap[chargingPileId] = chargingData
					}
				} else {
					// 数组为空，处理空数组的情况
				}

			}
		}

		// for rows.Next() {
		// 	var chargingPile models.ChargingPile
		// 	if err := rows.Scan(&chargingPile.Id, &chargingPile.Mqtt_id, &chargingPile.Energy); err != nil {
		// 		Logger.Error(err)
		// 	}
		// 	fmt.Printf("Id:%d, Mqtt_id:%d, Energy:%f\n", chargingPile.Id, chargingPile.Mqtt_id, chargingPile.Energy)
		// 	/*
		// 		收到電壓訊號 表示此台電表狀態是 standby 或是 charging (charging讓電流控制) 去和 OneMinDataMap的Status最後一筆去比
		// 		1.非 standby or charging => new status = standby
		// 		2.standby  => new status = standby
		// 		3.charging => new status = charging
		// 		4.
		// 	*/
		// }
		defer rows.Close()
	}
}

func CurrentDataProcessing(kind, topic, payload string) {
}

func EnergyDataProcessing(kind, topic, payload string) {
}
