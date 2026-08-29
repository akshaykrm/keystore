package httpx

import (
	"encoding/json"
	"fmt"
	"io"
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
	_ = json.NewEncoder(w).Encode(response)
}

func Write(w http.ResponseWriter, response SuccessResponse, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func Json(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Failed to send back message: %v", err)

	}

}

func Decode[T any](payload io.ReadCloser, target *T) error {
	if err := json.NewDecoder(payload).Decode(target); err != nil {
		return fmt.Errorf("invalid body")
	}
	return nil
}
