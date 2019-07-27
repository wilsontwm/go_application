package models

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"os"
	"github.com/joho/godotenv"
	"time"
	"log"
	"fmt"
)

var db *gorm.DB // database

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid)
}

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

	migrateDatabase()
}


// Datebase migration
func migrateDatabase() {
	db.Debug().AutoMigrate(
		&User{}, 
		&Company{},
		&Role{},
		&CompanyUser{},
	) 

	// Add foreign key
	db.Model(&Role{}).AddForeignKey("company_id", "companies(id)", "CASCADE", "RESTRICT")
	db.Model(&CompanyUser{}).AddForeignKey("company_id", "companies(id)", "CASCADE", "RESTRICT")
	db.Model(&CompanyUser{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
	db.Model(&CompanyUser{}).AddForeignKey("role_id", "roles(id)", "RESTRICT", "RESTRICT")

	// Add index

}

func GetDB() *gorm.DB {
	return db
}