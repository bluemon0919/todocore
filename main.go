package main

import (
	"log"

	"todocore/entity"
	"todocore/todo"
)

func main() {
	ent := entity.NewMap()
	if ent == nil {
		log.Fatal("file create error")
		return
	}
	srv := todo.NewServer(":8080", ent)
	go srv.StartService()

	client, _ := todo.NewClient("http://localhost:8080")

	//ui := NewMenu(client, os.Stdin)
	ui := NewGUI(client)
	//ui := NewHandler(client, ":8000")
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
