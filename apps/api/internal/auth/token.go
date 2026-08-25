package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	User string `json:"user"`
	jwt.RegisteredClaims
}

func CreateNewToken(claims Claims) (string, error) {
	key := []byte("super secret key")

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := t.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("Failed to create token: %w", err)
	}

	return s, nil
}

func ParseToken(token string) (*Claims, error) {

	t, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte("super secret key"), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Parsing token failed: %w", err)
	} else if claims, ok := t.Claims.(*Claims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("unknown claims type, cannot proceed")
	}

}
