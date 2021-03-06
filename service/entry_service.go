package service

import (
	"applichic.com/chic_secret/config"
	"applichic.com/chic_secret/model"
	uuid "github.com/satori/go.uuid"
	"time"
)

type EntryService struct {
}

// Save an entry
func (e *EntryService) Save(entry *model.Entry) error {
	err := config.DB.Save(&entry).Error
	return err
}

// GetEntry Get an entry
func (e *EntryService) GetEntry(entryId uuid.UUID) (model.Entry, error) {
	var entry model.Entry
	err := config.DB.
		Where("id = ?", entryId).
		Find(&entry).Error

	return entry, err
}

// GetEntriesToSynchronize Get the modified entries linked to the vault
func (e *EntryService) GetEntriesToSynchronize(userId uuid.UUID, lastSync time.Time) ([]model.Entry, error) {
	var entries []model.Entry
	err := config.DB.
		Joins("left join vaults on vaults.id = entries.vault_id").
		Joins("left join vault_users on vault_users.vault_id = vaults.id").
		Where("(vaults.user_id = ? or vault_users.user_id = ?) and entries.updated_at > ?", userId, userId, lastSync).
		Find(&entries).Error

	return entries, err
}

// GetEntriesFromVault Get the all entries linked to the vault
func (e *EntryService) GetEntriesFromVault(userId uuid.UUID) ([]model.Entry, error) {
	var entries []model.Entry
	err := config.DB.
		Joins("left join vaults on vaults.id = entries.vault_id").
		Joins("left join vault_users on vault_users.vault_id = vaults.id").
		Where("vaults.user_id = ? or vault_users.user_id = ?", userId, userId).
		Find(&entries).Error

	return entries, err
}

func (e *EntryService) DeleteFromUser(userId uuid.UUID) {
	config.DB.Exec("delete from entries "+
		"using vaults "+
		"where vaults.id = entries.vault_id "+
		"and vaults.user_id = ?", userId)
}
