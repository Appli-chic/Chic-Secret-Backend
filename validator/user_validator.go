package validator

import "applichic.com/chic_secret/model"

type SaveUserForm struct {
	User model.User `validate:"required"`
}
