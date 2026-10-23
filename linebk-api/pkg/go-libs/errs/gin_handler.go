package errs

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/dto"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/enums/api"
)

// GinHandleError standardizes and sends an HTTP error response based on the concrete error type.
func GinHandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	_ = c.Error(err)

	switch e := err.(type) {
	case ValidationError:
		c.JSON(e.StatusCode, dto.ApiResponse{
			Status:     e.Status,
			StatusCode: e.StatusCode,
			Message:    e.Message,
			Errors:     e.Errors,
		})
		return
	case AppError:
		c.JSON(e.StatusCode, dto.ApiResponse{
			Status:     e.Status,
			StatusCode: e.StatusCode,
			Message:    e.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, dto.ApiResponse{
		Status:     api.StatusUnexpected,
		StatusCode: http.StatusInternalServerError,
		Message:    "Unexpected Error",
	})
}
