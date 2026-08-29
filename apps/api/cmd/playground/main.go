package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/akshaykrm/keystore/apps/api/internal/user/workspace"
	_ "modernc.org/sqlite"
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
		fmt.Printf("Connecting to db failed: %v", err)
	}

	workspaceRepository := workspace.NewRepository(db)
	data := workspace.Workspace{
		Name:      "Test",
		Slug:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = workspaceRepository.Create(data)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Workspace inserted")

	workspaces, err := workspaceRepository.GetAll()

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("workspace list: ", len(workspaces))

	workspaceById, err := workspaceRepository.GetById("1")

	if err != nil {
		fmt.Printf("get by id failed :%v\n", workspaceById)
	}

	fmt.Println("workspace: ", workspaceById)
	updated := workspace.Workspace{
		ID:        "1",
		Name:      "Test updated",
		Slug:      "name updated",
		UpdatedAt: time.Now(),
	}

	err = workspaceRepository.UpdateById(updated)
	if err != nil {
		fmt.Printf("update failed: %v", err)
	} else {
		fmt.Println("updated workspace")
	}

	deleteItem := workspace.Workspace{
		ID: "1",
	}
	err = workspaceRepository.DeleteById(deleteItem)
	if err != nil {
		fmt.Printf("delete failed: %v", err)
	} else {
		fmt.Println("delete workspace")
	}

}
