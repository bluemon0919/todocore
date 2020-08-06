package main

import (
	"fmt"
	"log"
	"os"

	"todocore/entity"
	"todocore/todo"
)

func main() {
	projectID := os.Getenv("PROJECT_ID")
	if "" == projectID {
		fmt.Println("PROJECT_ID is not set")
		os.Exit(0)
	}
	ent := entity.NewDatastore(projectID, "RegistrationData")
	if ent == nil {
		log.Fatal("file create error")
		return
	}
	srv := todo.NewServer(":8080", ent)
	srv.StartService()
}
