package main

import (
	"errors"
)

// EntityMap データを管理する
type EntityMap struct {
	m  map[int]EntityItem
	id int
}

// NewEntityMap creates Entity
func NewEntityMap() *EntityMap {
	return &EntityMap{}
}

// NewID 新しいIDを返す
func (e *EntityMap) NewID() int {
	e.id++
	return e.id
}

// Add Entityにアイテムを追加する
func (e *EntityMap) Add(ei *EntityItem) error {
	if e.m == nil {
		e.m = make(map[int]EntityItem)
	}
	e.m[ei.Key] = *ei
	return nil
}

// Delete Entityから指定キーを削除する
func (e *EntityMap) Delete(key int) error {
	if _, ok := e.m[key]; !ok {
		return errors.New("key does not exist")
	}
	delete(e.m, key)
	return nil
}

// Update Entityの指定のキーを入力ステータスでアップデートする
func (e *EntityMap) Update(key, status int) error {
	if _, ok := e.m[key]; !ok {
		return errors.New("key does not exist")
	}
	item := e.m[key]
	item.Status = status
	e.m[key] = item
	return nil
}

// Get Entityからアイテムを取得する
func (e *EntityMap) Get(status int) (eis []EntityItem, err error) {
	key := func(p1, p2 *EntityItem) bool {
		return p1.Key < p2.Key
	}

	for _, ei := range e.m {
		switch status {
		case ACTIVE:
			if ei.Status != ACTIVE {
				continue
			}
		case COMPLETE:
			if ei.Status != COMPLETE {
				continue
			}
		}
		eis = append(eis, ei)
	}
	By(key).Sort(eis)
	return
}
