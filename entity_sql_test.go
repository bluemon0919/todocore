package main

import (
	"log"
	"os"
	"testing"
)

const TestEntitySQLTestFileName = "test.db"

func TestNewEntitySQL(t *testing.T) {

	expect := TestEntitySQLTestFileName
	// ファイルを削除する
	if _, err := os.Stat(expect); !os.IsNotExist(err) {
		if err = os.Remove(expect); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}

	ent := NewEntitySQL(expect)
	if ent.filename != expect {
		log.Fatalf("file name does not match. %s\n", ent.filename)
	}
	if ent.db == nil {
		log.Fatal("db is nil")
	}

	// ファイルを削除する
	if _, err := os.Stat(expect); !os.IsNotExist(err) {
		if err = os.Remove(expect); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}
}

func TestAdd(t *testing.T) {

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}

	ent := NewEntitySQL(TestEntitySQLTestFileName)

	item := EntityItem{
		Detail: "hogehoge",
		Status: ACTIVE,
	}
	if err := ent.Add(&item); err == nil {
		log.Fatal("register even if the Title is empty")
	}

	item = EntityItem{
		Title:  "hoge",
		Status: ACTIVE,
	}
	if err := ent.Add(&item); err == nil {
		log.Fatal("register even if the Detail is empty")
	}

	item = EntityItem{
		Title:  "hoge",
		Detail: "hogehoge",
	}
	if err := ent.Add(&item); err == nil {
		log.Fatal("register even if the Status is empty")
	}

	item = EntityItem{
		Title:  "hoge",
		Detail: "hogehoge",
		Status: ACTIVE,
	}
	if err := ent.Add(&item); err != nil {
		log.Fatal("registration failed")
	}

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}
}

func TestDelete(t *testing.T) {

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}

	ent := NewEntitySQL(TestEntitySQLTestFileName)

	item := EntityItem{
		Title:  "hoge",
		Detail: "hogehoge",
		Status: ACTIVE,
	}

	if err := ent.Add(&item); err != nil {
		log.Fatal("registration failed")
	}

	if err := ent.Delete(1); err != nil {
		log.Fatal("delete failed")
	}

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}
}

func TestUpdate(t *testing.T) {

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}

	ent := NewEntitySQL(TestEntitySQLTestFileName)

	item := EntityItem{
		Title:  "hoge",
		Detail: "hogehoge",
		Status: ACTIVE,
	}

	if err := ent.Add(&item); err != nil {
		log.Fatal("registration failed")
	}

	if err := ent.Update(1, COMPLETE); err != nil {
		log.Fatal("update failed")
	}

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}
}

func TestGet(t *testing.T) {

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}

	ent := NewEntitySQL(TestEntitySQLTestFileName)

	item := EntityItem{
		Title:  "hoge",
		Detail: "hogehoge",
		Status: ACTIVE,
	}

	if err := ent.Add(&item); err != nil {
		log.Fatal("registration failed")
	}

	items, err := ent.Get(ACTIVE)
	if err != nil {
		log.Fatal("update failed")
	}
	for _, item := range items {
		switch {
		case item.Title != "hoge":
			log.Fatal("title does not match")
		case item.Detail != "hogehoge":
			log.Fatal("detail does not match")
		case item.Status != ACTIVE:
			log.Fatal("status does not match")
		}
	}

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}
}
