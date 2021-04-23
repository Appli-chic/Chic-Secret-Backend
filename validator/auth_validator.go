package validator

type AskCodeForm struct {
	Email string `validate:"required,email"`
}

type LoginUserForm struct {
	Email string `validate:"required,email"`
	Token int    `validate:"required"`
}

type RefreshingTokenForm struct {
	RefreshToken string `validate:"required"`
}
