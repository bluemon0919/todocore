package main

import (
	"log"

	"todotool/entity"
)

func main() {
	ent := entity.NewEntitySQL("database.db")
	if ent == nil {
		log.Fatal("file create error")
		return
	}

	td := NewTODO(ent)

	handler := NewHandler(td, ":8080")
	if err := handler.Run(); err != nil {
		log.Fatal(err)
	}
}
