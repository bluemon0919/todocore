package entity

import (
	"log"
	"testing"
	"time"
)

var table = []Item{
	{0, "title1", "detail1", ACTIVE, time.Time{}, "ABC"},
	{1, "title2", "detail2", COMPLETE, time.Time{}, "ABC"},
	{2, "title3", "detail3", COMPLETE, time.Time{}, "ABC"},
	{3, "title4", "detail4", ACTIVE, time.Time{}, "ABC"},
	{4, "title5", "detail5", ACTIVE, time.Time{}, "ABC"},
}

func TestGet(t *testing.T) {
	entity := &EntityMap{}

	entity.m = make(map[int]Item)
	for i, ei := range table {
		entity.m[i] = ei
	}

	expect, err := entity.Get(ALL)
	if err != nil {
		log.Fatal(err)
	}
	for i, ei := range expect {
		if table[i].Key != ei.Key {
			log.Fatal(ei.Key)
		}
		if table[i].Title != ei.Title {
			log.Fatal(ei.Title)
		}
		if table[i].Detail != ei.Detail {
			log.Fatal(ei.Detail)
		}
		if table[i].Status != ei.Status {
			log.Fatal(ei.Status)
		}
	}
}
