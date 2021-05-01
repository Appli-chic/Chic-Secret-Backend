package validator

import "applichic.com/chic_secret/model"

type SaveTagForm struct {
	Tags []model.Tag `validate:"required"`
}
