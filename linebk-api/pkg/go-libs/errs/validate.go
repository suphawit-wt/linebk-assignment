package errs

import (
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/validator"
)

// GinBindAndValidate binds a JSON payload to a request struct and validates its input format.
func GinBindAndValidate(c *gin.Context, req any, vdt *validator.Instance) error {
	enTrans, _ := vdt.Uni.GetTranslator("en")

	if err := c.ShouldBindJSON(&req); err != nil {
		errorMessages := validator.TranslateError(err, enTrans, reflect.TypeOf(req).Elem())
		return NewValidationError("Invalid input format.", errorMessages)
	}

	return nil
}
