package main

import (
	"database/sql"
	"log"

	"github.com/akshaykrm/keystore/apps/api/internal/user"
	_ "modernc.org/sqlite"

	"fmt"
	"net/http"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./db/keystore.db")

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func main() {
	db, err := Connect()
	if err != nil {
		log.Fatalf("Connecting to db failed: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is live\n"))
	})

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userController := user.NewController(userService)

	user.RegisterRoutes(mux, userController)

	// authController := auth.NewController()
	// auth.RegisterRoutes(mux, authController)

	fmt.Println("Server started on port 3000")
	err = http.ListenAndServe(":3000", mux)

	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
}
