package service

import (
	"applichic.com/chic_secret/config"
	"applichic.com/chic_secret/model"
)

type TokenService struct {
}

// Save the token
func (t *TokenService) Save(token model.Token) (model.Token, error) {
	config.DB.NewRecord(token)
	err := config.DB.Create(&token).Error
	return token, err
}

// FetchTokenByUserId Fetch a token the user's email
func (t *TokenService) FetchTokenByUserId(userId interface{}) (model.Token, error) {
	token := model.Token{}
	err := config.DB.Where("user_id = ? AND is_valid = ?", userId, true).First(&token).Error
	return token, err
}
