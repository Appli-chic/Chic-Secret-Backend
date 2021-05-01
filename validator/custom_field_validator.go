package validator

import "applichic.com/chic_secret/model"

type SaveCustomFieldForm struct {
	CustomFields []model.CustomField `validate:"required"`
}
