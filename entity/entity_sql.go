package entity

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tenntenn/sqlite"
)

// EntitySQL データを管理する
type EntitySQL struct {
	filename string // TODO:ファイルの拡張子はシステム側でつけてもいいかも
	db       *sql.DB
	id       int
}

// NewSQL creates Entity
func NewSQL(filename string) *EntitySQL {
	db, err := sql.Open(sqlite.DriverName, filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	const sql = `CREATE TABLE IF NOT EXISTS item (
		key   INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		detail TEXT NOT NULL,
		status INTEGER NOT NULL);`

	if _, err = db.Exec(sql); err != nil {
		return nil
	}

	return &EntitySQL{
		filename: filename,
		db:       db,
	}
}

// NewID 新しいIDを返す
func (e *EntitySQL) NewID() int {
	e.id++
	return e.id
}

// Add Entityにアイテムを追加する
func (e *EntitySQL) Add(ei *EntityItem) error {
	switch {
	case ei.Title == "":
		return errors.New("Title is empty")
	case ei.Detail == "":
		return errors.New("Detail is empty")
	case ei.Status == 0:
		return errors.New("Status is empty")
	}
	const sql = "INSERT INTO item(title, detail, status) values (?,?,?)"
	_, err := e.db.Exec(sql, ei.Title, ei.Detail, ei.Status)
	return err
}

// Delete Entityから指定キーを削除する
func (e *EntitySQL) Delete(key int) error {
	sql := "DELETE FROM item WHERE key = ?"
	_, err := e.db.Exec(sql, key)
	return err
}

// Update Entityの指定のキーを入力ステータスでアップデートする
func (e *EntitySQL) Update(key, status int) error {
	sql := "UPDATE item SET status = ? WHERE key = ?"
	_, err := e.db.Exec(sql, status, key)
	return err
}

// Get Entityからアイテムを取得する
func (e *EntitySQL) Get(status int) (eis []EntityItem, err error) {
	// TODO : キーでソートする
	const sql = "SELECT * FROM item WHERE status = ?"
	rows, err := e.db.Query(sql, status)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var item EntityItem
		if err := rows.Scan(&item.Key, &item.Title, &item.Detail, &item.Status); err != nil {
			return nil, err
		}
		eis = append(eis, item)
	}
	return
}
