package validator

import "applichic.com/chic_secret/model"

type SaveVaultsForm struct {
	Vaults []model.Vault `validate:"required"`
}
