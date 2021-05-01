package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type CustomField struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Value     string    `gorm:"type:varchar(255);not null"`
	EntryID   uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (customField *CustomField) BeforeCreate(tx *gorm.DB) (err error) {
	if customField.ID == uuid.Nil {
		customField.ID = uuid.NewV4()
	}

	return
}
