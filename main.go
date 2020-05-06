package main

import (
	"log"
)

func main() {
	ent := NewEntity()
	td := NewTODO(ent)
	menu := NewGUI(td)
	if err := menu.Run(); err != nil {
		log.Fatal(err)
	}
}
