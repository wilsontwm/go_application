package models

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jinzhu/gorm"
	"os"
	"github.com/joho/godotenv"
	"log"
	"fmt"
)

var db *gorm.DB // database

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	dbUri := fmt.Sprintf("postgres://%v@%v:%v/%v?sslmode=disable&password=%v", username, dbHost, dbPort, dbName, password)
	
	// Making connection to the database
	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		log.Println(err)
	}

	db = conn
	db.Debug().AutoMigrate(
		&User{}, //&Contact{}
	) // Datebase migration
}

func GetDB() *gorm.DB {
	return db
}