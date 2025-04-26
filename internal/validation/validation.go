package validation

import (
	"encoding/json"
	"errors"
	"hackfest-uc/internal/domain/dto"
	"strings"

	"github.com/go-playground/validator"
)

type InputValidation struct {
	Validator *validator.Validate
}

func NewInputValidation() *InputValidation {
	return &InputValidation{
		Validator: validator.New(),
	}
}

func (v *InputValidation) Validate(data interface{}) error {
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
				errField.Message = err.Field() + " must be minimum " + err.Param() + " characters"
			case "required":
				errField.FieldName = strings.ToLower(err.Field())
				errField.Message = err.Field() + " cannot be blank"
			case "alpha":
				errField.FieldName = strings.ToLower(err.Field())
				errField.Message = err.Field() + " must contain only letters"
			}
			validationErrors = append(validationErrors, errField)
		}
	}
	if len(validationErrors) == 0 {
		return nil
	}
	marshaledErr, _ := json.Marshal(validationErrors)
	return errors.New(string(marshaledErr))
}
