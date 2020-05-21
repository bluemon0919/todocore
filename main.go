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

	handler := NewHandler(td, ":8080")
	if err := handler.Run(); err != nil {
		log.Fatal(err)
	}
}
