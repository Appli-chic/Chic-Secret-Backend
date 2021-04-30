package service

import (
	"applichic.com/chic_secret/config"
	"applichic.com/chic_secret/model"
	uuid "github.com/satori/go.uuid"
	"time"
)

type VaultService struct {
}

// Save a vault
func (v *VaultService) Save(vault *model.Vault) error {
	err := config.DB.Save(&vault).Error
	return err
}

// GetVaultsToSynchronize Get the modified vaults
func (v *VaultService) GetVaultsToSynchronize(userId uuid.UUID, lastSync time.Time) ([]model.Vault, error) {
	var vaults []model.Vault
	err := config.DB.Where("user_id = ? AND updated_at > ?", userId, lastSync).Find(&vaults).Error
	return vaults, err
}

// GetUserVaults Get the all vaults linked to the user
func (v *VaultService) GetUserVaults(userId uuid.UUID) ([]model.Vault, error) {
	var vaults []model.Vault
	err := config.DB.Where("user_id = ?", userId).Find(&vaults).Error
	return vaults, err
}
