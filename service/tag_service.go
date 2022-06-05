package service

import (
	"applichic.com/chic_secret/config"
	"applichic.com/chic_secret/model"
	uuid "github.com/satori/go.uuid"
	"time"
)

type TagService struct {
}

// Save a Tag
func (t *TagService) Save(tag *model.Tag) error {
	err := config.DB.Save(&tag).Error
	return err
}

// GetTag Get a tag
func (t *TagService) GetTag(tagId uuid.UUID) (model.Tag, error) {
	var tag model.Tag
	err := config.DB.
		Where("id = ?", tagId).
		Find(&tag).Error

	return tag, err
}

// GetTagsToSynchronize Get the modified tags linked to the vault
func (t *TagService) GetTagsToSynchronize(userId uuid.UUID, lastSync time.Time) ([]model.Tag, error) {
	var tags []model.Tag
	err := config.DB.
		Joins("left join vaults on vaults.id = tags.vault_id").
		Joins("left join vault_users on vault_users.vault_id = vaults.id").
		Where("(vaults.user_id = ? or vault_users.user_id = ?) AND tags.updated_at > ?", userId, userId, lastSync).
		Find(&tags).Error

	return tags, err
}

// GetTagsFromVault Get the all the tags linked to the vault
func (t *TagService) GetTagsFromVault(userId uuid.UUID) ([]model.Tag, error) {
	var tags []model.Tag
	err := config.DB.
		Joins("left join vaults on vaults.id = tags.vault_id").
		Joins("left join vault_users on vault_users.vault_id = vaults.id").
		Where("vaults.user_id = ? or vault_users.user_id = ?", userId, userId).
		Find(&tags).Error

	return tags, err
}

func (t *TagService) DeleteFromUser(userId uuid.UUID) {
	config.DB.Exec("delete from tags "+
		"using vaults "+
		"where vaults.id = tags.vault_id "+
		"and vaults.user_id = ?", userId)
}
