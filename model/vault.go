package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Vault struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;"`
	Name       string     `gorm:"type:varchar(255);not null"`
	Signature  string     `gorm:"type:varchar(255);not null"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null"`
	Categories []Category `gorm:"foreignKey:VaultID"`
	Entries    []Entry    `gorm:"foreignKey:VaultID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
}

func (vault *Vault) BeforeCreate(tx *gorm.DB) (err error) {
	if vault.ID == uuid.Nil {
		vault.ID = uuid.NewV4()
	}

	return
}
