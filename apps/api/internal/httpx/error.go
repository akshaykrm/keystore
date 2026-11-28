package httpx

// Response is the common response wrapper for all http responses.
type ResponseOut[T any] struct {
	// Message is a short human readable summary of the result.
	Message string `json:"message,omitempty"`

	// Data hold the successful results, shapes depend on the endpoint.
	Data T `json:"data,omitempty"`

	// Error holds the failure response, if the request failed.
	Error string `json:"error,omitempty"`

	// Status is the HTTP status code, it is not sent in the JSON body.
	Status int `json:"-"`
}

// ErrorResponse is the body returned for any failed response
type ErrorResponse struct {

	// Message is a short human readable summary of the result.
	Message string `json:"message,omitempty"`

	// Error holds the failure response, if the request failed.
	Error string `json:"error,omitempty"`
}

type Response[T any] struct {
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Error   error  `json:"error,omitempty"`
	Status  int    `json:"-"`
}

func (r Response[T]) toOut() ResponseOut[T] {
	var msg string

	if r.Error != nil {
		msg = r.Error.Error()
	}

	return ResponseOut[T]{
		Message: r.Message,
		Data:    r.Data,
		Error:   msg,
		Status:  r.Status,
	}

}
