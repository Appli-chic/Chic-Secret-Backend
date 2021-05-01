package service

import (
	"applichic.com/chic_secret/config"
	"applichic.com/chic_secret/model"
)

type UserService struct {
}

// FetchUserById Fetch a user from it's ID
func (u *UserService) FetchUserById(userId interface{}) (model.User, error) {
	user := model.User{}
	err := config.DB.Select("id, email").Where("id = ?", userId).First(&user).Error
	return user, err
}

// FetchUserByEmail Fetch a user from it's email
func (u *UserService) FetchUserByEmail(email string) (model.User, error) {
	user := model.User{}
	err := config.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

// FetchUserFromRefreshToken Fetch a user from the refresh token linked to this account
func (u *UserService) FetchUserFromRefreshToken(refreshToken string) (model.User, error) {
	user := model.User{}
	err := config.DB.
		Joins("left join tokens on tokens.user_id = users.id").
		Where("tokens.token = ?", refreshToken).
		First(&user).Error
	return user, err
}

// Save a user
func (u *UserService) Save(user *model.User) error {
	err := config.DB.Create(&user).Error
	return err
}
