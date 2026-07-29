package main

import (
	"database/sql"
	"log"

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
	_, err := Connect()

	if err != nil {
		log.Fatalf("Connecting to db failed: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is live\n"))
	})

	fmt.Println("Server starting...")
	http.ListenAndServe(":5000", mux)

}
