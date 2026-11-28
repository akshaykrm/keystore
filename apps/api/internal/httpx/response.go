package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Json[T any](w http.ResponseWriter, response Response[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	if err := json.NewEncoder(w).Encode(response.toOut()); err != nil {
		fmt.Printf("Failed to send back message: %v", err)
	}

}

func Decode[T any](payload io.ReadCloser, target *T) error {
	if err := json.NewDecoder(payload).Decode(target); err != nil {
		return fmt.Errorf("invalid body")
	}
	return nil
}
