package validation

import (
	"hackfest-uc/internal/domain/dto"
	"strings"

	"github.com/go-playground/validator"
)

type InputValidation struct {
	Validator *validator.Validate
}

func NewInputValidation() *InputValidation {
	validate := validator.New()
	return &InputValidation{
		Validator: validate,
	}
}

func (v *InputValidation) Validate(data interface{}) []dto.ErrorInputResponse {
	var validationErrors []dto.ErrorInputResponse

	err := v.Validator.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var errField dto.ErrorInputResponse
			switch err.Tag() {
			case "email":
				errField.FieldName = strings.ToLower(err.Field())
				errField.Message = "Email format is invalid"
			case "min":
				errField.FieldName = strings.ToLower(err.Field())
				errField.Message = err.Field() + " must be at least " + err.Param() + " characters"
			case "required":
				errField.FieldName = strings.ToLower(err.Field())
				errField.Message = err.Field() + " is required"
			case "gt":
				errField.FieldName = strings.ToLower(err.Field())
				errField.Message = err.Field() + " must be greater than " + err.Param()
			case "oneof":
				errField.FieldName = strings.ToLower(err.Field())
				errField.Message = err.Field() + " must be one of: " + strings.ReplaceAll(err.Param(), " ", ", ")
			default:
				errField.FieldName = strings.ToLower(err.Field())
				errField.Message = "Invalid value for " + err.Field()
			}
			validationErrors = append(validationErrors, errField)
		}
	}

	return validationErrors
}
