package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/akshaykrm/keystore/apps/api/db"
	"github.com/akshaykrm/keystore/apps/api/internal/membership"
	"github.com/akshaykrm/keystore/apps/api/internal/products"
	"github.com/akshaykrm/keystore/apps/api/internal/user"
	"github.com/akshaykrm/keystore/apps/api/internal/workspace"

	_ "modernc.org/sqlite"
)

func Connect() (*sql.DB, error) {
	conn, err := sql.Open("sqlite", "./db/keystore.db")

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func main() {
	conn, err := Connect()
	if err != nil {
		log.Fatalf("Connecting to db failed: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is live\n"))
	})

	queries := db.New(conn)

	productService := products.NewService(queries)
	productController := products.NewController(productService)

	membershipRepo := membership.NewRepository(conn)
	membershipService := membership.NewService(membershipRepo)

	workspaceRepo := workspace.NewRepository(conn)
	workspaceService := workspace.NewService(workspaceRepo, membershipService)
	workspaceController := workspace.NewController(workspaceService)

	userRepo := user.NewRepository(conn)
	userService := user.NewService(userRepo, workspaceService, membershipService)
	userController := user.NewController(userService)

	// Register Roues
	user.RegisterRoutes(mux, userController)
	workspace.RegisterRoutes(mux, workspaceController)
	products.RegisterRoutes(mux, productController)

	fmt.Println("Server started on port 3000")
	err = http.ListenAndServe(":3000", mux)

	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
}
