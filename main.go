package main

import "log"

func main() {
	ent := NewEntitySQL("database.db")
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
