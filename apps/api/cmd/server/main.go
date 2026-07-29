package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is live\n"))
	})

	fmt.Println("Server starting...")
	http.ListenAndServe(":5000", mux)

}
