package service

import (
	"applichic.com/chic_secret/config"
	"applichic.com/chic_secret/model"
	uuid "github.com/satori/go.uuid"
	"time"
)

type CustomFieldService struct {
}

// Save a CustomField
func (c *CustomFieldService) Save(customField *model.CustomField) error {
	err := config.DB.Save(&customField).Error
	return err
}

// GetCustomFieldsToSynchronize Get the modified custom fields linked to the vault
func (c *CustomFieldService) GetCustomFieldsToSynchronize(userId uuid.UUID, lastSync time.Time) ([]model.CustomField, error) {
	var customFields []model.CustomField
	err := config.DB.
		Joins("left join entries on entries.id = custom_fields.entry_id").
		Joins("left join vaults on vaults.id = entries.vault_id").
		Where("vaults.user_id = ? AND custom_fields.updated_at > ?", userId, lastSync).
		Find(&customFields).Error

	return customFields, err
}

// GetCustomFieldsFromVault Get the all the custom fields linked to the vault
func (c *CustomFieldService) GetCustomFieldsFromVault(userId uuid.UUID) ([]model.CustomField, error) {
	var customFields []model.CustomField
	err := config.DB.
		Joins("left join entries on entries.id = custom_fields.entry_id").
		Joins("left join vaults on vaults.id = entries.vault_id").
		Where("vaults.user_id = ?", userId).
		Find(&customFields).Error

	return customFields, err
}
