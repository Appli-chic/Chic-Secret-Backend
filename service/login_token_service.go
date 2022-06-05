package service

import (
	"applichic.com/chic_secret/config"
	"applichic.com/chic_secret/model"
	uuid "github.com/satori/go.uuid"
	"time"
)

type LoginTokenService struct {
}

// Save the login token in the database
func (l *LoginTokenService) Save(loginToken model.LoginToken) (model.LoginToken, error) {
	err := config.DB.Create(&loginToken).Error
	return loginToken, err
}

// DeleteAllForUser Delete all the login tokens linked to the user
func (l *LoginTokenService) DeleteAllForUser(userId uuid.UUID) error {
	err := config.DB.Where("id = ", userId).Delete(&model.LoginToken{}).Error
	return err
}

// FetchTokenNotExpiredByUserId Fetch the login token that are not expired for this user
func (l *LoginTokenService) FetchTokenNotExpiredByUserId(userId uuid.UUID, code int) (model.LoginToken, error) {
	loginToken := model.LoginToken{}
	err := config.DB.Where("user_id = ? AND token = ? AND expire_at >= ?", userId, code, time.Now()).First(&loginToken).Error
	return loginToken, err
}

func (l *LoginTokenService) DeleteFromUser(userId uuid.UUID) {
	config.DB.Exec("delete from login_tokens "+
		"where login_tokens.user_id = ?", userId)
}
