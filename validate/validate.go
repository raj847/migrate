package validate

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Validate(v interface{}) error {
	return validate.Struct(v)
}
