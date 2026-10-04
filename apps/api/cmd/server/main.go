package main

import (
	"database/sql"
	"log"

	"github.com/akshaykrm/keystore/apps/api/internal/membership"
	"github.com/akshaykrm/keystore/apps/api/internal/user"
	"github.com/akshaykrm/keystore/apps/api/internal/workspace"
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

	workspaceRepo := workspace.NewRepository(db)
	workspaceService := workspace.NewService(workspaceRepo)
	workspaceController := workspace.NewController(workspaceService)

	membershipRepo := membership.NewRepository(db)
	membershipService := membership.NewService(membershipRepo)

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, workspaceService, membershipService)
	userController := user.NewController(userService)

	// Register Roues
	user.RegisterRoutes(mux, userController)
	workspace.RegisterRoutes(mux, workspaceController)

	fmt.Println("Server started on port 3000")
	err = http.ListenAndServe(":3000", mux)

	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
}
