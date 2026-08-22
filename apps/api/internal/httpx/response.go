package httpx

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
}

func Error2(w http.ResponseWriter, response ErrorResponse, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func Write(w http.ResponseWriter, response SuccessResponse, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
