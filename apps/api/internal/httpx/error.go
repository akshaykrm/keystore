package httpx

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   error  `json:"error,omitempty"`
	Status  int    `json:"-"`
}

func (r Response) toOut() ResponseOut {
	var msg string

	if r.Error != nil {
		msg = r.Error.Error()
	}

	return ResponseOut{
		Message: r.Message,
		Data:    r.Data,
		Error:   msg,
		Status:  r.Status,
	}

}

type ResponseOut struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Status  int    `json:"-"`
}
