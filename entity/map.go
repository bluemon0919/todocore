package entity

import (
	"errors"
	"time"
)

// EntityMap データを管理する
type EntityMap struct {
	m  map[int]Item
	id int
}

// NewMap creates Entity
func NewMap() *EntityMap {
	return &EntityMap{}
}

// NewID 新しいIDを返す
func (ent *EntityMap) NewID() int {
	ent.id++
	return ent.id
}

// Add Entityにアイテムを追加する
func (ent *EntityMap) Add(item *Item) error {
	if ent.m == nil {
		ent.m = make(map[int]Item)
	}
	ent.m[item.Key] = *item
	return nil
}

// Delete Entityから指定キーを削除する
func (ent *EntityMap) Delete(key int) error {
	if _, ok := ent.m[key]; !ok {
		return errors.New("key does not exist")
	}
	delete(ent.m, key)
	return nil
}

// Update Entityの指定のキーを入力ステータスでアップデートする
func (ent *EntityMap) Update(key, status int) error {
	if _, ok := ent.m[key]; !ok {
		return errors.New("key does not exist")
	}
	item := ent.m[key]
	item.Status = status
	ent.m[key] = item
	return nil
}

// Get Entityからアイテムを取得する
func (ent *EntityMap) Get(status int) (items []Item, err error) {
	key := func(p1, p2 *Item) bool {
		return p1.Key < p2.Key
	}

	for _, ei := range ent.m {
		if status != ei.Status {
			continue
		}
		items = append(items, ei)
	}
	By(key).Sort(items)
	return
}

// GetDate 期間を指定してアイテムを取得する
func (ent *EntityMap) GetDate(start, end time.Time) (items []Item, err error) {
	err = nil
	for _, data := range ent.m {
		if end.After(data.Date) && start.Before(data.Date) {
			items = append(items, data)
		}
	}

	key := func(p1, p2 *Item) bool {
		return p1.Key < p2.Key
	}
	By(key).Sort(items)
	return
}

// GetDate 期間を指定してアイテムを取得する
func (ent *EntityMap) GetAfterDate(start time.Time) (items []Item, err error) {
	err = nil
	for _, data := range ent.m {
		if start.Before(data.Date) {
			items = append(items, data)
		}
	}

	key := func(p1, p2 *Item) bool {
		return p1.Key < p2.Key
	}
	By(key).Sort(items)
	return
}
