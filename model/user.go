package model

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID          uuid.UUID    `gorm:"type:uuid;primary_key;"`
	Email       string       `gorm:"type:varchar(255);unique_index"`
	Tokens      []Token      `gorm:"foreignkey:UserID"`
	LoginTokens []LoginToken `gorm:"foreignkey:UserID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid)
}
