package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type LoginToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Token     int
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	ExpireAt  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (loginToken *LoginToken) BeforeCreate(tx *gorm.DB) (err error) {
	if loginToken.ID == uuid.Nil {
		loginToken.ID = uuid.NewV4()
	}

	return
}
