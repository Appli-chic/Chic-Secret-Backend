package model

import (
	"gorm.io/gorm"
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID           uuid.UUID    `gorm:"type:uuid;primary_key;"`
	Email        string       `gorm:"type:varchar(255);unique_index"`
	Tokens       []Token      `gorm:"foreignKey:UserID"`
	LoginTokens  []LoginToken `gorm:"foreignKey:UserID"`
	Vaults       []Vault      `gorm:"foreignKey:UserID"`
	IsSubscribed bool         `gorm:"not null"`
	Subscription string       `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.NewV4()
	}

	return
}
