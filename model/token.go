package model

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Token struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Token     string    `gorm:"type:varchar(36);unique_index"`
	IsValid   bool      `gorm:"not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (token *Token) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid)
}
