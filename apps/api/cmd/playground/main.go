package main

import (
	"fmt"

	"github.com/akshaykrm/keystore/apps/api/internal/auth"
)

func main() {
	password := "test"
	hashed, err := auth.HashPassword(password)

	if err != nil {
		fmt.Printf("hashing failed %v", err)
		return
	}
	fmt.Printf("Hashed Password: %s\n", hashed)

	err = auth.CompareHashedPassword(hashed, "test2")
	fmt.Println(err)
}
