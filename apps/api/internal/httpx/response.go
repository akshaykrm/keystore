package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
