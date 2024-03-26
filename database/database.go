package models

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() *sql.DB {
	connStr := "user=alan dbname=side sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
