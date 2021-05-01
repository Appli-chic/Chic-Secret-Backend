package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Tag struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;"`
	Name      string     `gorm:"type:varchar(255);not null"`
	VaultID   uuid.UUID  `gorm:"type:uuid;not null"`
	EntryTags []EntryTag `gorm:"foreignKey:TagID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (tag *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	if tag.ID == uuid.Nil {
		tag.ID = uuid.NewV4()
	}

	return
}
