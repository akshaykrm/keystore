package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/akshaykrm/keystore/apps/api/internal/httpx"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func GetClaims(r *http.Request) (*Claims, bool) {
	claims, ok := r.Context().Value(ClaimsKey).(*Claims)
	return claims, ok

}

func IsAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		unauthorised := func() {
			res := httpx.ErrorResponse{
				Message: "Restricted from access",
				Error:   AccessDenied.Error(),
			}
			httpx.Error2(w, res, http.StatusUnauthorized)
		}

		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			unauthorised()
			return
		}
		prefix, token, found := strings.Cut(bearerToken, " ")
		if !found || token == "" || !strings.EqualFold(prefix, "Bearer") {
			unauthorised()
			return
		}

		claims, err := ParseToken(token)
		if err != nil {
			unauthorised()
			fmt.Printf("Parse Token Failed %v\n", err)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)

		next(w, r.WithContext(ctx))
	}
}
