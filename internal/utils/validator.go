package utils

import "github.com/go-playground/validator/v10"

type Validator struct {
	v *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{v: validator.New()}
}

func (cv *Validator) Validate(i any) error {
	return cv.v.Struct(i)
}
