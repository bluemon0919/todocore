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

	td := todo.NewTODO(ent)

	//ui := NewMenu(td, os.Stdin)
	//ui := NewGUI(td)
	ui := NewHandler(td, ":8080")
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
