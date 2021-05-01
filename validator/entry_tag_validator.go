package validator

import "applichic.com/chic_secret/model"

type SaveEntryTagForm struct {
	EntryTags []model.EntryTag `validate:"required"`
}
