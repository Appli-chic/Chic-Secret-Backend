package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Icon      int       `gorm:"not null"`
	Color     string    `gorm:"type:varchar(255);not null"`
	IsTrash   bool      `gorm:"not null"`
	VaultID   uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (category *Category) BeforeCreate(scope *gorm.Scope) error {
	if category.ID == uuid.Nil {
		uuid := uuid.NewV4()
		return scope.SetColumn("ID", uuid)
	}

	return nil
}
