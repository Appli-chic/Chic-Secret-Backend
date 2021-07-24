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

// GetEntryTag Get an entry category
func (e *EntryTagService) GetEntryTag(entryId uuid.UUID, tagId uuid.UUID) (model.EntryTag, error) {
	var entryTag model.EntryTag
	err := config.DB.
		Where("entry_id = ? AND tag_id = ?", entryId, tagId).
		Find(&entryTag).Error

	return entryTag, err
}

// GetEntryTagsToSynchronize Get the modified entry tags linked to the vault
func (e *EntryTagService) GetEntryTagsToSynchronize(userId uuid.UUID, lastSync time.Time) ([]model.EntryTag, error) {
	var entryTags []model.EntryTag
	err := config.DB.
		Joins("left join entries on entries.id = entry_tags.entry_id").
		Joins("left join vaults on vaults.id = entries.vault_id").
		Joins("left join vault_users on vault_users.vault_id = vaults.id").
		Where("(vaults.user_id = ? or vault_users.user_id = ?) AND entry_tags.updated_at > ?", userId, userId, lastSync).
		Find(&entryTags).Error

	return entryTags, err
}

// GetEntryTagsFromVault Get the all the entry tags linked to the vault
func (e *EntryTagService) GetEntryTagsFromVault(userId uuid.UUID) ([]model.EntryTag, error) {
	var entryTags []model.EntryTag
	err := config.DB.
		Joins("left join entries on entries.id = entry_tags.entry_id").
		Joins("left join vaults on vaults.id = entries.vault_id").
		Joins("left join vault_users on vault_users.vault_id = vaults.id").
		Where("vaults.user_id = ? or vault_users.user_id = ?", userId, userId).
		Find(&entryTags).Error

	return entryTags, err
}
