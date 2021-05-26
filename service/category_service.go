package service

import (
	"applichic.com/chic_secret/config"
	"applichic.com/chic_secret/model"
	uuid "github.com/satori/go.uuid"
	"time"
)

type CategoryService struct {
}

// Save a category
func (c *CategoryService) Save(category *model.Category) error {
	err := config.DB.Save(&category).Error
	return err
}

// GetCategory Get a category
func (c *CategoryService) GetCategory(categoryId uuid.UUID) (model.Category, error) {
	var category model.Category
	err := config.DB.
		Where("id = ?", categoryId).
		Find(&category).Error

	return category, err
}

// GetCategoriesToSynchronize Get the modified categories linked to the vault
func (c *CategoryService) GetCategoriesToSynchronize(userId uuid.UUID, lastSync time.Time) ([]model.Category, error) {
	var categories []model.Category
	err := config.DB.
		Joins("left join vaults on vaults.id = categories.vault_id").
		Where("vaults.user_id = ? AND categories.updated_at > ?", userId, lastSync).
		Find(&categories).Error

	return categories, err
}

// GetCategoriesFromVault Get the all categories linked to the vault
func (c *CategoryService) GetCategoriesFromVault(userId uuid.UUID) ([]model.Category, error) {
	var categories []model.Category
	err := config.DB.
		Joins("left join vaults on vaults.id = categories.vault_id").
		Where("vaults.user_id = ?", userId).
		Find(&categories).Error

	return categories, err
}
