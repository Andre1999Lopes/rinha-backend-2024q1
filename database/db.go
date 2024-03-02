package database

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func CreateConnection() {
	var conn string = "host=db user=admin password=123 dbname=rinha port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(conn))

	for err != nil {
		time.Sleep(2 * time.Second)
		DB, err = gorm.Open(postgres.Open(conn))
	}
}
