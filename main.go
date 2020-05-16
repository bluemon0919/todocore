package main

import "log"

func main() {
	ent := NewEntitySQL("database.db")
	//ent := NewEntityJSON("tmp.json")
	if ent == nil {
		log.Fatal("file create error")
		return
	}
	//ent := NewEntityMap()
	td := NewTODO(ent)
	menu := NewGUI(td)
	if err := menu.Run(); err != nil {
		log.Fatal(err)
	}
}
