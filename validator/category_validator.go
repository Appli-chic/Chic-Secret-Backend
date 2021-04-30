package validator

import "applichic.com/chic_secret/model"

type SaveCategoriesForm struct {
	Categories []model.Category `validate:"required"`
}
