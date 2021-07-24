package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type VaultUser struct {
	VaultID   uuid.UUID `gorm:"type:uuid;primaryKey;"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
