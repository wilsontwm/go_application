package models

import (
	//"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type Role struct {
	Base
	Name string
	IsAdmin bool `gorm:"default:false"`
	CompanyID uuid.UUID `gorm:"type:uuid;not null;"`
}
