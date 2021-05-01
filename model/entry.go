package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Entry struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name       string    `gorm:"type:varchar(255);not null"`
	Username   string    `gorm:"type:varchar(255);not null"`
	Hash       string    `gorm:"type:varchar(255);not null"`
	Comment    string    `gorm:"not null"`
	VaultID    uuid.UUID `gorm:"type:uuid;not null"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
}

func (entry *Entry) BeforeCreate(tx *gorm.DB) (err error) {
	if entry.ID == uuid.Nil {
		entry.ID = uuid.NewV4()
	}

	return
}
