package main

import (
	"fmt"
	"time"

	"github.com/akshaykrm/keystore/apps/api/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	fmt.Println("Token generation checking: ")

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(24 * time.Hour)
	claims := auth.Claims{
		User:      "test@gmail.com",
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	token, _ := auth.CreateNewToken(claims)
	fmt.Println("Gen Token: ", token)
	parsedClaims, err := auth.ParseToken(token)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Claims: ", parsedClaims.ExpiresAt, parsedClaims.IssuedAt, parsedClaims.User)

	fmt.Println("Token generation checking done!!! ")

}
