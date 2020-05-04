package main

import (
	"log"
	"os"
)

func main() {
	ent := NewEntity()
	td := NewTODO(ent)
	menu := NewMenu(td, os.Stdin)
	if err := menu.Run(); err != nil {
		log.Fatal(err)
	}
}
