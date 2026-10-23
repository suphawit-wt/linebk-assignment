package card

// DebitCardResponse return a debit card item
type DebitCardResponse struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Issuer      string `json:"issuer"`
	Color       string `json:"color"`
	BorderColor string `json:"borderColor,omitzero"`
	Number      string `json:"number,omitzero"`
}

// FromListModel is for map []DebitCardInfo model to []DebitCardResponse response
func (i *DebitCardResponse) FromListModel(model []DebitCardInfo) []DebitCardResponse {
	responses := make([]DebitCardResponse, len(model))

	for i := range model {
		responses[i] = DebitCardResponse{
			Name:        model[i].Name.String,
			Status:      model[i].Status.String,
			Issuer:      model[i].Issuer.String,
			Color:       model[i].Color.String,
			BorderColor: model[i].BorderColor.String,
			Number:      model[i].Number.String,
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
