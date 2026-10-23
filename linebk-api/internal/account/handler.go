package account

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/enums/api"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/errs"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/validator"
)

// Handler defines HTTP request handler methods for account.
type Handler interface {
	GetAccountsByUserID(c *gin.Context)
}

type handler struct {
	vdt    *validator.Instance
	accSvc Service
}

// NewHandler creates a new instance of Handler with the provided dependencies.
func NewHandler(
	vdt *validator.Instance,
	accSvc Service,
) Handler {
	return &handler{vdt: vdt, accSvc: accSvc}
}

// GetAccountsByUserID godoc
//
//	@Summary		Get accounts by user id
//	@Description	Get a list of accounts by user id
//	@Tags			Accounts
//	@Security		Bearer
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ApiResponseWithData{data=[]InfoResponse}
//	@Failure		500	{object}	ApiResponse
//	@Router			/v1/accounts [get]
func (h *handler) GetAccountsByUserID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := c.MustGet("userID").(string)
	if !ok {
		errs.GinHandleError(c, errs.NewUnexpectedError())
		return
	}

	accountsInfo, err := h.accSvc.GetAccountsByUserID(ctx, userID)
	if err != nil {
		errs.GinHandleError(c, err)
		return
	}

	var accountsInfoRes InfoResponse

	res := ApiResponseWithData[[]InfoResponse]{
		Status:     api.StatusSuccess,
		StatusCode: http.StatusOK,
		Data:       accountsInfoRes.FromListModel(accountsInfo),
	}

	c.JSON(http.StatusOK, res)
}
