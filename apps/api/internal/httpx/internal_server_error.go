package httpx

import (
	"log"
	"net/http"

	"github.com/akshaykrm/keystore/apps/api/internal/apperr"
)

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("ERROR: %s %s: %v\n", r.Method, r.URL.Path, err)
	Json(w, Response{
		Message: "an unexpected error occurred. Please try again later",
		Error:   apperr.ErrInternalServerError,
		Status:  http.StatusInternalServerError,
	})
}
