package main

import (
	"log"
	"testing"
)

var table = []EntityItem{
	{0, "title1", "detail1", ACTIVE},
	{1, "title2", "detail2", COMPLETE},
	{2, "title3", "detail3", COMPLETE},
	{3, "title4", "detail4", ACTIVE},
	{4, "title5", "detail5", ACTIVE},
}

func TestGet(t *testing.T) {
	entity := &Entity{}

	entity.m = make(map[int]EntityItem)
	for i, ei := range table {
		entity.m[i] = ei
	}

	expect, err := entity.Get(ALL)
	if err != nil {
		log.Fatal(err)
	}
	for i, ei := range expect {
		if table[i].key != ei.key {
			log.Fatal(ei.key)
		}
		if table[i].title != ei.title {
			log.Fatal(ei.title)
		}
		if table[i].detail != ei.detail {
			log.Fatal(ei.detail)
		}
		if table[i].status != ei.status {
			log.Fatal(ei.status)
		}
	}
}
