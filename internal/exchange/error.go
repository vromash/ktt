package exchange

func NewErrorResponse(err string) ErrorResponse {
	return ErrorResponse{Error: err}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
