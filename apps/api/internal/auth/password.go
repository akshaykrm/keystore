package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("hash failed: %w", err)
	}

	return string(hashedPassword), nil
}

func CompareHashedPassword(hashedPassword, password string) bool {
	bytePassword := []byte(password)
	byteHash := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePassword)
	return err != nil
}
