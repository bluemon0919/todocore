package main

import (
	"io/ioutil"

	"github.com/rivo/tview"
)

// DetailPanel パネルへのプレビューを行う
type DetailPanel struct {
	*tview.TextView
}

// NewDetailPanel creates DetailPanel
func NewDetailPanel() *DetailPanel {
	p := &DetailPanel{
		TextView: tview.NewTextView(),
	}

	p.SetBorder(true).
		SetTitle("preview").
		SetTitleAlign(tview.AlignLeft)

	return p
}

// UpdateView 表示する
func (p *DetailPanel) UpdateView(name string) {
	var content string
	b, err := ioutil.ReadFile(name)
	if err != nil {
		content = err.Error()
	} else {
		content = string(b)
	}

	p.Clear().SetText(content)
}
