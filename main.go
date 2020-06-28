package main

import (
	"fmt"
	"log"
	"os"

	"todocore/entity"
	"todocore/todo"
)

func main() {
	//ent := entity.NewMap()
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
	go srv.StartService()

	client, _ := todo.NewClient("http://localhost:8080")

	ui := NewMenu(client, os.Stdin)
	//ui := NewGUI(client)
	//ui := NewHandler(client, ":8000")
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
