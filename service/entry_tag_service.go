package service

import (
	"applichic.com/chic_secret/config"
	"applichic.com/chic_secret/model"
	uuid "github.com/satori/go.uuid"
	"time"
)

type EntryTagService struct {
}

// Save a Tag
func (e *EntryTagService) Save(entryTag *model.EntryTag) error {
	err := config.DB.Save(&entryTag).Error
	return err
}

// GetEntryTagsToSynchronize Get the modified entry tags linked to the vault
func (e *EntryTagService) GetEntryTagsToSynchronize(userId uuid.UUID, lastSync time.Time) ([]model.EntryTag, error) {
	var entryTags []model.EntryTag
	err := config.DB.
		Joins("left join entries on entries.id = entry_tags.entry_id").
		Joins("left join vaults on vaults.id = entries.vault_id").
		Where("vaults.user_id = ? AND entry_tags.updated_at > ?", userId, lastSync).
		Find(&entryTags).Error

	return entryTags, err
}

// GetEntryTagsFromVault Get the all the entry tags linked to the vault
func (e *EntryTagService) GetEntryTagsFromVault(userId uuid.UUID) ([]model.EntryTag, error) {
	var entryTags []model.EntryTag
	err := config.DB.
		Joins("left join entries on entries.id = entry_tags.entry_id").
		Joins("left join vaults on vaults.id = entries.vault_id").
		Where("vaults.user_id = ? ", userId).
		Find(&entryTags).Error

	return entryTags, err
}
