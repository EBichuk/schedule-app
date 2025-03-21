package validators

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	Validator *validator.Validate
}

func (val *CustomValidator) Validate(i interface{}) error {
	return val.Validator.Struct(i)
}
