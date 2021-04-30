package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Vault struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name      string    `gorm:"type:varchar(255)"`
	Signature string    `gorm:"type:varchar(255)"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (vault *Vault) BeforeCreate(scope *gorm.Scope) error {
	if vault.ID == uuid.Nil {
		uuid := uuid.NewV4()
		return scope.SetColumn("ID", uuid)
	}

	return nil
}
