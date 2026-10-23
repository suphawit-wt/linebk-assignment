package banner

// Response return an banner item
type Response struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

// FromListModel is for map []Banner model to []BannerResponse response
func (i *Response) FromListModel(model []Banner) []Response {
	responses := make([]Response, len(model))

	for i := range model {
		responses[i] = Response{
			Title:       model[i].Title.String,
			Description: model[i].Description.String,
			Image:       model[i].Image.String,
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
