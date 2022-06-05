package service

import (
	"applichic.com/chic_secret/config"
	"applichic.com/chic_secret/model"
	uuid "github.com/satori/go.uuid"
	"time"
)

type VaultUserService struct {
}

// Save a vault user
func (v *VaultUserService) Save(vaultUser *model.VaultUser) error {
	err := config.DB.Save(&vaultUser).Error
	return err
}

// GetVaultUser Get a vault user
func (v *VaultUserService) GetVaultUser(vaultId uuid.UUID, userId uuid.UUID) (model.VaultUser, error) {
	var vaultUser model.VaultUser
	err := config.DB.
		Joins("left join vaults on vaults.id = vault_users.vault_id").
		Where("vault_id = ? AND (vaults.user_id = ? or vault_users.user_id = ?)", vaultId, userId, userId).
		Find(&vaultUser).Error

	return vaultUser, err
}

// GetVaultUsersToSynchronize Get the modified vault users
func (v *VaultUserService) GetVaultUsersToSynchronize(userId uuid.UUID, lastSync time.Time) ([]model.VaultUser, error) {
	var vaultUsers []model.VaultUser
	err := config.DB.
		Joins("left join vaults on vaults.id = vault_users.vault_id").
		Joins("left join users on users.id = vault_users.user_id").
		Where("(vaults.user_id = ? or vault_users.user_id = ?) AND vault_users.updated_at > ?", userId, userId, lastSync).
		Find(&vaultUsers).Error

	return vaultUsers, err
}

// GetVaultUsers Get the all the vault users
func (v *VaultUserService) GetVaultUsers(userId uuid.UUID) ([]model.VaultUser, error) {
	var vaultUsers []model.VaultUser
	err := config.DB.
		Joins("left join vaults on vaults.id = vault_users.vault_id").
		Joins("left join users on users.id = vault_users.user_id").
		Where("vaults.user_id = ? or vault_users.user_id = ? ", userId, userId).
		Find(&vaultUsers).Error

	return vaultUsers, err
}

func (v *VaultUserService) DeleteFromUser(userId uuid.UUID) {
	config.DB.Exec("delete from vault_users "+
		"using vaults "+
		"where vaults.id = vault_users.vault_id "+
		"and vaults.user_id = ?", userId)
}
