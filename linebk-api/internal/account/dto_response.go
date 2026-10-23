package account

import (
	"strings"
)

// InfoResponse return information of account
type InfoResponse struct {
	Type          string   `json:"type"`
	Amount        string   `json:"amount"`
	Currency      string   `json:"currency"`
	AccountNumber string   `json:"accountNumber"`
	Issuer        string   `json:"issuer"`
	Progress      *uint8   `json:"progress,omitempty"`
	Flags         []string `json:"flags,omitempty"`
	Color         string   `json:"color"`
	IsMainAccount bool     `json:"isMainAccount"`
}

// FromListModel is for map []Info model to []InfoResponse response
func (i *InfoResponse) FromListModel(model []Info) []InfoResponse {
	responses := make([]InfoResponse, len(model))

	for i := range model {
		responses[i] = InfoResponse{
			Type:          model[i].Type,
			Amount:        model[i].Amount.StringFixed(2),
			Currency:      model[i].Currency,
			AccountNumber: model[i].AccountNumber,
			Issuer:        model[i].Issuer,
			Progress:      model[i].Progress,
			Color:         model[i].Color.String,
			IsMainAccount: model[i].IsMainAccount,
		}

		if model[i].RawFlags.Valid && model[i].RawFlags.String != "" {
			responses[i].Flags = strings.Split(model[i].RawFlags.String, ",")
		}
	}

	return responses
}

// ApiResponse represents the standard structure for API responses, including status, status code, message, and optional error details.
type ApiResponse struct {
	Status     string        `json:"status"`
	StatusCode int           `json:"statusCode"`
	Message    string        `json:"message"`
	Errors     []ErrResponse `json:"errors,omitempty" swaggerignore:"true"`
}

// ApiResponseWithData represents the standard structure for API responses that include data, along with status and status code.
type ApiResponseWithData[T any] struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
	Data       T      `json:"data"       swaggerignore:"true"`
}

// ApiResponseWithDataAndPagination includes both data and pagination metadata for API responses
type ApiResponseWithDataAndPagination[T any] struct {
	Data       T               `json:"data"                 swaggerignore:"true"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
	Status     string          `json:"status"`
	StatusCode int             `json:"statusCode"`
}

// PaginationMeta provides metadata for paginated responses, indicating the current page, items per page, and whether there are next or previous pages.
type PaginationMeta struct {
	Page       int  `json:"page,omitempty"`
	PerPage    int  `json:"perPage,omitempty"`
	TotalCount int  `json:"totalCount"`
	HasNext    bool `json:"hasNext"`
	HasPrev    bool `json:"hasPrev,omitempty"`
}

// ErrResponse represents the structure for individual error responses within the API.
type ErrResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
