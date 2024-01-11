package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	return &Validator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		errv, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		if len(errv) > 0 {
			err := errv[0]
			switch err.Tag() {
			case "required":
				return fmt.Errorf("%v: is required", err.Field())
			default:
				return err
			}
		}
	}
	return nil
}
