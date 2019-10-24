package db

import (
	"os"

	"github.com/jinzhu/gorm"
)

// Connect is connect at DB by Gorm
func Connect() *gorm.DB {
	DBMS := os.Getenv("DBMS")
	USERNAME := os.Getenv("USERNAME")
	PW := os.Getenv("PW")
	// PROTOCOL := ""
	DBNAME := os.Getenv("DBNAME")
	CONNECT := USERNAME + ":" + PW + "@/" + DBNAME
	println(USERNAME, CONNECT)
	db, err := gorm.Open(DBMS, CONNECT)
	db.LogMode(true)
	if err != nil {
		panic(err.Error())
	}
	return db
}
