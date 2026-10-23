package transaction

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/enums/api"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/errs"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/validator"
)

// Handler defines HTTP request handler methods for transaction.
type Handler interface {
	GetRecentTransactions(c *gin.Context)
}

type handler struct {
	vdt   *validator.Instance
	txSvc Service
}

// NewHandler creates a new instance of Handler with the provided dependencies.
func NewHandler(
	vdt *validator.Instance,
	txSvc Service,
) Handler {
	return &handler{vdt: vdt, txSvc: txSvc}
}

// GetRecentTransactions godoc
//
//	@Summary		Get recent transactions by user id
//	@Description	Get a list of recent transactions by user id
//	@Tags			Transactions
//	@Security		Bearer
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ApiResponseWithData{data=[]RecentTransactionResponse}
//	@Failure		500	{object}	ApiResponse
//	@Router			/v1/transactions/recent [get]
func (h *handler) GetRecentTransactions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := c.MustGet("userID").(string)
	if !ok {
		errs.GinHandleError(c, errs.NewUnexpectedError())
		return
	}

	transactions, err := h.txSvc.GetRecentTransactions(ctx, userID)
	if err != nil {
		errs.GinHandleError(c, err)
		return
	}

	var recentTransactions RecentTransactionResponse

	res := ApiResponseWithData[[]RecentTransactionResponse]{
		Status:     api.StatusSuccess,
		StatusCode: http.StatusOK,
		Data:       recentTransactions.FromListModel(transactions),
	}

	c.JSON(http.StatusOK, res)
}
