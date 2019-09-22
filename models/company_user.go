package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type CompanyUser struct {
	CompanyID   uuid.UUID  `gorm:"type:uuid;not null;primary_key"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;primary_key"`
	RoleID      uuid.UUID  `gorm:"type:uuid"`
	LastVisited *time.Time `gorm:"index:last_visited"`
}
