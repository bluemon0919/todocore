package entity

import (
	"sort"
)

// Entity operates entity
type Entity interface {
	NewID() int
	Add(ei *EntityItem) error
	Delete(key int) error
	Update(key, status int) error
	Get(status int) (eis []EntityItem, err error)
}

// EntityItem Entityに書き込むアイテム
type EntityItem struct {
	Key    int
	Title  string
	Detail string
	Status int
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
