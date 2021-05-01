package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type EntryTag struct {
	EntryID   uuid.UUID `gorm:"type:uuid;primaryKey;"`
	TagID     uuid.UUID `gorm:"type:uuid;primaryKey;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
