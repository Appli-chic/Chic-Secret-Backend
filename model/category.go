package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Icon      int       `gorm:"not null"`
	Color     string    `gorm:"type:varchar(255);not null"`
	IsTrash   bool      `gorm:"not null"`
	VaultID   uuid.UUID `gorm:"type:uuid;not null"`
	Entry     Entry     `gorm:"foreignKey:CategoryID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (category *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if category.ID == uuid.Nil {
		category.ID = uuid.NewV4()
	}

	return
}
