package service

func VoltageDataProcessing(kind, topic, payload string) {
	// rows, err := db.Query("SELECT * FROM users")
	// if err != nil {
	// 	panic(err)
	// }
	// for rows.Next() {
	// 	var chargingPile models.ChargingPile
	// 	if err := rows.Scan(&chargingPile.ID, &chargingPile); err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Printf("ID: %d, Name: %s, Age: %s\n", user.ID, user.Username, user.Password)
	// }
	// defer rows.Close()
}

func CurrentDataProcessing(kind, topic, payload string) {
}

func EnergyDataProcessing(kind, topic, payload string) {
}
