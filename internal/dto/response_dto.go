package dto

type (
	ErrorResponse struct {
		Message string `json:"messsage"`
		Details any    `json:"details,omitempty"`
	}

	DetailsError struct {
		Title   string `json:"field"`
		Message string `json:"message"`
	}

	Response struct {
		Success bool           `json:"success"`
		Error   *ErrorResponse `json:"error,omitempty"`
		Data    any            `json:"data,omitempty"`
	}

	PageSort struct {
		Columns  []string `json:"columns"`
		OrderDir string   `json:"OrderDir"`
	}

	PageFilter struct {
		Field string `json:"field"`
		Value any    `json:"value"`
	}

	PageInfo struct {
		Page      int          `json:"page"`
		Limit     int          `json:"limit"`
		OrderBy   []string     `json:"order_by"`
		OrderSort string       `json:"order_sort"`
		Filters   []PageFilter `json:"filters"`
		TotalRow  int          `json:"total_row"`
	}

	PaginatadResponse[T any] struct {
		Entries  []T      `json:"entries"`
		PageInfo PageInfo `json:"page_info"`
	}

	PaginationQuery struct {
		Limit string `form:"limit,omitempty" binding:"omitempty,number"`
		Page  string `form:"page,omitempty" binding:"omitempty,number"`
	}
)

func ResponseError(errRes ErrorResponse) Response {
	return Response{
		Success: false,
		Error:   &errRes,
	}
}
