package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	User      string           `json:"user"`
	ID        int64            `json:"id"`
	IssuedAt  *jwt.NumericDate `json:"issued_at"`
	ExpiresAt *jwt.NumericDate `json:"expires_at"`
	jwt.RegisteredClaims
}

// TODO: the key needs to come from someplace and not be part of code
func CreateNewToken(ID int64, Email string) (string, error) {
	key := []byte("super secret key")
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(24 * time.Hour)
	claims := Claims{
		User:      Email,
		ID:        ID,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

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
