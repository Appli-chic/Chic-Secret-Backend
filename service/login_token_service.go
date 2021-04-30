package service

import (
	"applichic.com/chic_secret/config"
	"applichic.com/chic_secret/model"
	uuid "github.com/satori/go.uuid"
	"time"
)

type LoginTokenService struct {
}

func (l *LoginTokenService) Save(loginToken model.LoginToken) (model.LoginToken, error) {
	config.DB.NewRecord(loginToken)
	err := config.DB.Create(&loginToken).Error
	return loginToken, err
}

func (l *LoginTokenService) DeleteAllForUser(userId uuid.UUID) error {
	err := config.DB.Unscoped().Delete(&model.LoginToken{}).Where("id = ", userId).Error
	return err
}

func (l *LoginTokenService) FetchTokenNotExpiredByUserId(userId uuid.UUID, code int) (model.LoginToken, error) {
	loginToken := model.LoginToken{}
	err := config.DB.Where("user_id = ? AND token = ? AND expire_at >= ?", userId, code, time.Now()).First(&loginToken).Error
	return loginToken, err
}
