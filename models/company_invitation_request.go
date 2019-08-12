package models

import (
	"github.com/satori/go.uuid"
)

type CompanyInvitationRequest struct {
	Base
	CompanyID uuid.UUID `gorm:"type:uuid;not null;primary_key"`
	Email string `gorm:"not null;primary_key"`
	UserID *uuid.UUID `gorm:"type:uuid"`
	Status int `gorm:"default:'0'"`
}
