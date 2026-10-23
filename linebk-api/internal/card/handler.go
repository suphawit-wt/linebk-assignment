package card

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/enums/api"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/errs"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/validator"
)

// Handler defines HTTP request handler methods for card.
type Handler interface {
	GetDebitCards(c *gin.Context)
}

type handler struct {
	vdt     *validator.Instance
	cardSvc Service
}

// NewHandler creates a new instance of Handler with the provided dependencies.
func NewHandler(
	vdt *validator.Instance,
	cardSvc Service,
) Handler {
	return &handler{vdt: vdt, cardSvc: cardSvc}
}

// GetDebitCards godoc
//
//	@Summary		Get debit cards by user id
//	@Description	Get a list of debit cards by user id
//	@Tags			Debit Cards
//	@Security		Bearer
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ApiResponseWithData{data=[]DebitCardResponse}
//	@Failure		500	{object}	ApiResponse
//	@Router			/v1/debit-cards [get]
func (h *handler) GetDebitCards(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := c.MustGet("userID").(string)
	if !ok {
		errs.GinHandleError(c, errs.NewUnexpectedError())
		return
	}

	debitCards, err := h.cardSvc.GetDebitCards(ctx, userID)
	if err != nil {
		errs.GinHandleError(c, err)
		return
	}

	var debitCardsResponse DebitCardResponse

	res := ApiResponseWithData[[]DebitCardResponse]{
		Status:     api.StatusSuccess,
		StatusCode: http.StatusOK,
		Data:       debitCardsResponse.FromListModel(debitCards),
	}

	c.JSON(http.StatusOK, res)
}
