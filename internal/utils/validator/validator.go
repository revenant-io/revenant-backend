package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(data interface{}) error {
	err := validate.Struct(data)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return fmt.Errorf("validation error: %s", validationErrors[0].Error())
		}
		return fmt.Errorf("validation error: %v", err)
	}
	return nil
}
