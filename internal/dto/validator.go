package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	return nil
}
