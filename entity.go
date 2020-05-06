package main

import (
	"errors"
	"sort"
)

// Entity データを管理する
type Entity struct {
	m  map[int]EntityItem
	id int
}

var mm map[int]EntityItem

// EntityItem Entityに書き込むアイテム
type EntityItem struct {
	key    int
	title  string
	detail string
	status int
}

// NewEntity creates entity
func NewEntity() *Entity {
	return &Entity{}
}

// NewID 新しいIDを返す
func (e *Entity) NewID() int {
	e.id++
	return e.id
}

// Add Entityにアイテムを追加する
func (e *Entity) Add(ei *EntityItem) error {
	if e.m == nil {
		e.m = make(map[int]EntityItem)
	}
	e.m[ei.key] = *ei
	return nil
}

// Delete Entityから指定キーを削除する
func (e *Entity) Delete(key int) error {
	if _, ok := e.m[key]; !ok {
		return errors.New("key does not exist")
	}
	delete(e.m, key)
	return nil
}

// Update Entityの指定のキーを入力ステータスでアップデートする
func (e *Entity) Update(key, status int) error {
	if _, ok := e.m[key]; !ok {
		return errors.New("key does not exist")
	}
	item := e.m[key]
	item.status = status
	e.m[key] = item
	return nil
}

// Get Entityからアイテムを取得する
func (e *Entity) Get(status int) (eis []EntityItem, err error) {
	key := func(p1, p2 *EntityItem) bool {
		return p1.key < p2.key
	}

	for _, ei := range e.m {
		switch status {
		case ACTIVE:
			if ei.status != ACTIVE {
				continue
			}
		case COMPLETE:
			if ei.status != COMPLETE {
				continue
			}
		}
		eis = append(eis, ei)
	}
	By(key).Sort(eis)
	return
}

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(p1, p2 *EntityItem) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(items []EntityItem) {
	es := &entityItemSorter{
		items: items,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(es)
}

// planetSorter joins a By function and a slice of Planets to be sorted.
type entityItemSorter struct {
	items []EntityItem
	by    func(p1, p2 *EntityItem) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *entityItemSorter) Len() int {
	return len(s.items)
}

// Swap is part of sort.Interface.
func (s *entityItemSorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *entityItemSorter) Less(i, j int) bool {
	return s.by(&s.items[i], &s.items[j])
}
