package main

import (
	"database/sql"
	"fmt"

	"github.com/akshaykrm/keystore/apps/api/db"
	"github.com/akshaykrm/keystore/apps/api/internal/products"
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
		fmt.Fatalf("Connecting to conn failed: %v", err)

	}

	queries := db.New(conn)

	p := products.NewService(queries)

	fmt.Println("Creating Product")
	product, err := p.Create()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Product created %s", product.Name)
	}

	// p.GetAll()
	// p.UpdateByID(1)
	// p.GetOneByID(1)

}
