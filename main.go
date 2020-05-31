package main

import (
	"log"

	"todotool/entity"
	"todotool/todo"
)

func main() {
	ent := entity.NewEntitySQL("database.db")
	if ent == nil {
		log.Fatal("file create error")
		return
	}
	srv := todo.NewServer(":8080", ent)
	go srv.StartService()

	client, _ := todo.NewClient("http://localhost:8080")

	//ui := NewMenu(client, os.Stdin)
	//ui := NewGUI(client)
	ui := NewHandler(client, ":8000")
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
