package validator

import (
	"errors"
	"reflect"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/dto"
)

// TranslateError converts validation errors into a list of formatted DTO error responses using field JSON tag names.
func TranslateError(err error, trans ut.Translator, structType reflect.Type) []dto.ErrResponse {
	var validatorErrs validator.ValidationErrors
	if !errors.As(err, &validatorErrs) {
		return nil
	}

	errorMessages := make([]dto.ErrResponse, 0, len(validatorErrs))

	for _, e := range validatorErrs {
		jsonFieldName := getJSONFieldName(structType, e.Field())

		errorMessage := dto.ErrResponse{
			Field:   jsonFieldName,
			Message: e.Translate(trans),
		}
		errorMessages = append(errorMessages, errorMessage)
	}

	return errorMessages
}

func getJSONFieldName(structType reflect.Type, fieldName string) string {
	field, found := structType.FieldByName(fieldName)
	if !found {
		return fieldName
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName
	}
	return jsonTag
}
