package validator

import "applichic.com/chic_secret/model"

type SaveVaultUserForm struct {
	VaultUsers []model.VaultUser `validate:"required"`
}
