package entity

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// EntityJSON データを管理する
type EntityJSON struct {
	filename string
	m        map[int]EntityItem
	id       int
}

// Data 外部出力するJSONデータの構造体
type Data struct {
	ItemMap map[int]EntityItem
	ID      int
}

// NewEntityJSON creates entity
func NewEntityJSON(filename string) (ent *EntityJSON) {
	ent = &EntityJSON{
		filename: filename,
		id:       0,
	}
	if err := ent.Read(); err != nil {
		log.Fatal(err)
		ent = nil
	}
	return
}

// Read EntityItemをファイルから読み込みます
func (e *EntityJSON) Read() error {
	if _, err := os.Stat(e.filename); err != nil {
		fp, err := os.Create(e.filename)
		if err != nil {
			return nil
		}
		defer fp.Close()
		return nil
	}

	bytes, err := ioutil.ReadFile(e.filename)
	if err != nil {
		return err
	}

	if len(bytes) == 0 {
		return nil
	}

	var data Data
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	e.m = data.ItemMap
	e.id = data.ID
	return nil
}

// Write EntityItemをファイルに書き出します
func (e *EntityJSON) Write() error {
	data := Data{
		ItemMap: e.m,
		ID:      e.id,
	}

	bytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(e.filename, bytes, os.FileMode(0666))
}

// NewID 新しいIDを返す
func (e *EntityJSON) NewID() int {
	e.id++
	return e.id
}

// Add Entityにアイテムを追加する
func (e *EntityJSON) Add(ei *EntityItem) error {
	if e.m == nil {
		e.m = make(map[int]EntityItem)
	}
	e.m[ei.Key] = *ei
	return e.Write()
}

// Delete Entityから指定キーを削除する
func (e *EntityJSON) Delete(key int) error {
	if _, ok := e.m[key]; !ok {
		return errors.New("key does not exist")
	}
	delete(e.m, key)
	return e.Write()
}

// Update Entityの指定のキーを入力ステータスでアップデートする
func (e *EntityJSON) Update(key, status int) error {
	if _, ok := e.m[key]; !ok {
		return errors.New("key does not exist")
	}
	item := e.m[key]
	item.Status = status
	e.m[key] = item
	return e.Write()
}

// Get Entityからアイテムを取得する
func (e *EntityJSON) Get(status int) (eis []EntityItem, err error) {
	key := func(p1, p2 *EntityItem) bool {
		return p1.Key < p2.Key
	}

	for _, ei := range e.m {
		if status != ei.Status {
			continue
		}
		eis = append(eis, ei)
	}
	By(key).Sort(eis)
	return
}
