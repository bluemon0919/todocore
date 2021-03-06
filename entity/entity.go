package entity

import (
	"sort"
	"time"
)

// Entity operates entity
type Entity interface {
	Add(item *Item) error
	Delete(id int64) error
	Update(id int64, status int) error
	Get(status int) (items []Item, ids []int64, err error)
	GetDate(start, end time.Time) (items []Item, ids []int64, err error)
	GetAfterDate(base time.Time) (items []Item, ids []int64, err error)
	GetBeforeDate(base time.Time) (items []Item, ids []int64, err error)
}

// Item Entityに書き込むアイテム
type Item struct {
	Title     string
	Detail    string
	Status    int
	StartTime time.Time
	EndTime   time.Time
	StationID string
}

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(p1, p2 *Item) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(items []Item) {
	es := &entityItemSorter{
		items: items,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(es)
}

// planetSorter joins a By function and a slice of Planets to be sorted.
type entityItemSorter struct {
	items []Item
	by    func(p1, p2 *Item) bool // Closure used in the Less method.
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
