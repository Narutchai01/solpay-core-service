package utils

import "github.com/go-playground/validator/v10"

// utils/validator.go

// ฟังก์ชันสำหรับดึง Msg สั้นๆ
func FormatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// เอาแค่ error แรกที่เจอก็พอ
		err := validationErrors[0]
		switch err.Tag() {
		case "required":
			return err.Field() + " is required"
		case "len":
			return err.Field() + " must be " + err.Param() + " characters long"
			// ... case อื่นๆ
		}
		return err.Error() // default
	}
	return "Invalid input data"
}
