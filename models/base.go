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
var username, password, dbName, dbHost, dbPort string

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

	username = os.Getenv("db_user")
	password = os.Getenv("db_pass")
	dbName = os.Getenv("db_name")
	dbHost = os.Getenv("db_host")
	dbPort = os.Getenv("db_port")

	migrateDatabase()
}


// Datebase migration
func migrateDatabase() {
	db := GetDB()

	db.Debug().AutoMigrate(
		&User{}, 
		&Company{},
		&Role{},
		&CompanyUser{},
		&CompanyInvitationRequest{},
	) 

	// Migration scripts
	db.Model(&Role{}).AddForeignKey("company_id", "companies(id)", "CASCADE", "RESTRICT")
	db.Model(&CompanyUser{}).AddForeignKey("company_id", "companies(id)", "CASCADE", "RESTRICT")
	db.Model(&CompanyUser{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
	db.Model(&CompanyUser{}).AddForeignKey("role_id", "roles(id)", "RESTRICT", "RESTRICT")
	db.Model(&CompanyInvitationRequest{}).AddForeignKey("company_id", "companies(id)", "CASCADE", "RESTRICT")
	db.Model(&CompanyInvitationRequest{}).AddForeignKey("user_id", "users(id)", "SET NULL", "RESTRICT")
	db.Model(&User{}).DropColumn("birthday_string")
	db.Model(&User{}).DropColumn("token")
}

func GetDB() *gorm.DB {
	dbUri := fmt.Sprintf("postgres://%v@%v:%v/%v?sslmode=disable&password=%v", username, dbHost, dbPort, dbName, password)
	
	// Making connection to the database
	db, err := gorm.Open("postgres", dbUri)
	if err != nil {
		log.Println(err)
	}

	return db
}