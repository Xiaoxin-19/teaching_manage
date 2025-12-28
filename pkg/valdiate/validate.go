package valdiatex

import "github.com/go-playground/validator/v10"

var g_validator = validator.New()

func ValidateStruct(s interface{}) error {
	return g_validator.Struct(s)
}
