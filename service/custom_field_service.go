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

// GetCustomField Get a custom field
func (c *CustomFieldService) GetCustomField(CustomFieldId uuid.UUID) (model.CustomField, error) {
	var customField model.CustomField
	err := config.DB.
		Where("id = ?", CustomFieldId).
		Find(&customField).Error

	return customField, err
}

// GetCustomFieldsToSynchronize Get the modified custom fields linked to the vault
func (c *CustomFieldService) GetCustomFieldsToSynchronize(userId uuid.UUID, lastSync time.Time) ([]model.CustomField, error) {
	var customFields []model.CustomField
	err := config.DB.
		Joins("left join entries on entries.id = custom_fields.entry_id").
		Joins("left join vaults on vaults.id = entries.vault_id").
		Joins("left join vault_users on vault_users.vault_id = vaults.id").
		Where("(vaults.user_id = ? or vault_users.user_id = ?) AND custom_fields.updated_at > ?", userId, userId, lastSync).
		Find(&customFields).Error

	return customFields, err
}

// GetCustomFieldsFromVault Get the all the custom fields linked to the vault
func (c *CustomFieldService) GetCustomFieldsFromVault(userId uuid.UUID) ([]model.CustomField, error) {
	var customFields []model.CustomField
	err := config.DB.
		Joins("left join entries on entries.id = custom_fields.entry_id").
		Joins("left join vaults on vaults.id = entries.vault_id").
		Joins("left join vault_users on vault_users.vault_id = vaults.id").
		Where("vaults.user_id = ? or vault_users.user_id = ?", userId, userId).
		Find(&customFields).Error

	return customFields, err
}

func (c *CustomFieldService) DeleteFromUser(userId uuid.UUID) {
	config.DB.Exec("delete from custom_fields "+
		"using entries, vaults "+
		"where entries.id = custom_fields.entry_id "+
		"and vaults.id = entries.vault_id "+
		"and vaults.user_id = ?", userId)
}
