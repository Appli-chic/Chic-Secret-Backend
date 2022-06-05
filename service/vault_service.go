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

// GetVault Get a vault
func (v *VaultService) GetVault(vaultId uuid.UUID) (model.Vault, error) {
	var vault model.Vault
	err := config.DB.
		Where("id = ?", vaultId).
		Find(&vault).Error

	return vault, err
}

// GetVaultsToSynchronize Get the modified vaults
func (v *VaultService) GetVaultsToSynchronize(userId uuid.UUID, lastSync time.Time) ([]model.Vault, error) {
	var vaults []model.Vault
	err := config.DB.
		Joins("left join vault_users on vault_users.vault_id = vaults.id").
		Where("(vaults.user_id = ? or vault_users.user_id = ?) AND vaults.updated_at > ?", userId, userId, lastSync).Find(&vaults).Error
	return vaults, err
}

// GetUserVaults Get the all vaults linked to the user
func (v *VaultService) GetUserVaults(userId uuid.UUID) ([]model.Vault, error) {
	var vaults []model.Vault
	err := config.DB.
		Joins("left join vault_users on vault_users.vault_id = vaults.id").
		Where("vaults.user_id = ? or vault_users.user_id = ?", userId, userId).Find(&vaults).Error
	return vaults, err
}

func (v *VaultService) DeleteFromUser(userId uuid.UUID) {
	config.DB.Exec("delete from vaults "+
		"where vaults.user_id = ?", userId)
}
