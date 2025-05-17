package dto

type ErrorResponse struct {
	Message string `json:"messsage"`
	Details any    `json:"details,omitempty"`
}

type DetailsError struct {
	Title   string `json:"field"`
	Message string `json:"message"`
}

type Response struct {
	Success bool           `json:"success"`
	Error   *ErrorResponse `json:"error,omitempty"`
	Data    any            `json:"data,omitempty"`
}

func ResponseError(errRes ErrorResponse) Response {
	return Response{
		Success: false,
		Error:   &errRes,
	}
}
