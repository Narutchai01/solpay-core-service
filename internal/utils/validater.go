package utils

import "github.com/go-playground/validator/v10"

// FormatValidationError extracts a human-readable message from a validation error.
func FormatValidationError(err error) string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return "Invalid input data"
	}

	fe := validationErrors[0]
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "len":
		return fe.Field() + " must be " + fe.Param() + " characters long"
	default:
		return fe.Error()
	}
}
