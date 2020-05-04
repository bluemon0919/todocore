package main

import (
	"fmt"
	"io"
	"log"
	"testing"
)

type Climenutest struct {
	s []string
	i int
}

var count int

func (c *Climenutest) Read(p []byte) (n int, err error) {
	if count >= len(c.s) {
		return 0, io.EOF
	}
	fmt.Println(c.s[0])
	//	n = copy(p, []byte(c.s[c.i]))
	c.i++
	return 0, io.EOF
}

func TestNewIssue(t *testing.T) {
	//menu := NewMenu(nil, strings.NewReader("hoge"))
	c := &Climenutest{
		s: []string{"hoge", "fuga"},
		i: 0,
	}
	menu := NewMenu(nil, c)

	expect := menu.NewIssue()
	if expect == nil {
		log.Fatal(expect)
	}
}
