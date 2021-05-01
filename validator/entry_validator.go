package validator

import "applichic.com/chic_secret/model"

type SaveEntriesForm struct {
	Entries []model.Entry `validate:"required"`
}
