package banner

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/enums/api"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/errs"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/validator"
)

// Handler defines HTTP request handler methods for bannner.
type Handler interface {
	GetBanners(c *gin.Context)
}

type handler struct {
	vdt       *validator.Instance
	bannerSvc Service
}

// NewHandler creates a new instance of Handler with the provided dependencies.
func NewHandler(
	vdt *validator.Instance,
	bannerSvc Service,
) Handler {
	return &handler{vdt: vdt, bannerSvc: bannerSvc}
}

// GetBanners godoc
//
//	@Summary		Get banners by user id
//	@Description	Get a list of banners by user id
//	@Tags			Banners
//	@Security		Bearer
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ApiResponseWithData{data=[]Response}
//	@Failure		500	{object}	ApiResponse
//	@Router			/v1/banners [get]
func (h *handler) GetBanners(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := c.MustGet("userID").(string)
	if !ok {
		errs.GinHandleError(c, errs.NewUnexpectedError())
		return
	}

	banners, err := h.bannerSvc.GetBanners(ctx, userID)
	if err != nil {
		errs.GinHandleError(c, err)
		return
	}

	var bannersResponse Response

	res := ApiResponseWithData[[]Response]{
		Status:     api.StatusSuccess,
		StatusCode: http.StatusOK,
		Data:       bannersResponse.FromListModel(banners),
	}

	c.JSON(http.StatusOK, res)
}
