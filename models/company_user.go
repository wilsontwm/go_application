package models

import (
	//"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type CompanyUser struct {
	CompanyID uuid.UUID `gorm:"type:uuid;not null;primary_key"`
	UserID uuid.UUID `gorm:"type:uuid;not null;primary_key"`
	RoleID uuid.UUID `gorm:"type:uuid"`
}