package workspace

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func validate[T any](payload T) ValidationErrors {
	v := newValidator()
	err := v.Struct(payload)

	if err == nil {
		return nil
	}

	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return ValidationErrors{
			"-": "validation error",
		}
	}
	out := ValidationErrors{}
	for _, v := range errs {
		switch v.Tag() {
		case "required":
			out[v.Field()] = "is required"
		case "min":
			out[v.Field()] = fmt.Sprintf("must be at least %s characters", v.Param())
		}
	}
	return out
}
