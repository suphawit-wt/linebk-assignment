package user

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/enums/api"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/errs"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/validator"
)

// Handler defines HTTP request handler methods for user.
type Handler interface {
	GetUserInfo(c *gin.Context)
}

type handler struct {
	vdt     *validator.Instance
	userSvc Service
}

// NewHandler creates a new instance of Handler with the provided dependencies.
func NewHandler(
	vdt *validator.Instance,
	userSvc Service,
) Handler {
	return &handler{vdt: vdt, userSvc: userSvc}
}

// GetUserInfo godoc
//
//	@Summary		Get user info by id
//	@Description	Get a user info by id
//	@Tags			Users
//	@Security		Bearer
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ApiResponseWithData{data=InfoResponse}
//	@Failure		404	{object}	ApiResponse
//	@Failure		500	{object}	ApiResponse
//	@Router			/v1/user [get]
func (h *handler) GetUserInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := c.MustGet("userID").(string)
	if !ok {
		errs.GinHandleError(c, errs.NewUnexpectedError())
		return
	}

	userInfo, err := h.userSvc.GetUserInfo(ctx, userID)
	if err != nil {
		errs.GinHandleError(c, err)
		return
	}

	var infoRes InfoResponse

	res := ApiResponseWithData[*InfoResponse]{
		Status:     api.StatusSuccess,
		StatusCode: http.StatusOK,
		Data:       infoRes.FromModel(userInfo),
	}

	c.JSON(http.StatusOK, res)
}
