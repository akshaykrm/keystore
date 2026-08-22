package httpx

type ErrorResponse struct {
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
