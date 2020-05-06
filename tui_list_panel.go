package main

import (
	"fmt"

	"github.com/rivo/tview"
)

// ListPanel ファイル操作するための情報
type ListPanel struct {
	items []Item
	*tview.Table
}

// NewListPanel creates ListPanel
func NewListPanel() *ListPanel {
	p := &ListPanel{
		Table: tview.NewTable(),
	}

	p.SetBorder(true).
		SetTitle("TODO List").
		SetTitleAlign(tview.AlignLeft)

	p.SetSelectable(true, false)

	return p
}

// SetItems set files info
func (l *ListPanel) SetItems(items []Item) {
	l.items = items
}

// SelectedItem 現在選択しているファイル情報を取得する
func (l *ListPanel) SelectedItem() *Item {
	row, _ := l.Table.GetSelection()
	if row > len(l.items)-1 || row < 0 {
		return nil
	}
	return &l.items[row]
}

// Keybinding ListPanelのキーバインドを設定する
func (l *ListPanel) Keybinding(g *GUI) {

	// 子ウィンドウを定義する
	modal := tview.NewModal().
		SetText("Whitch operation?").
		AddButtons([]string{"Complete", "Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			item := l.SelectedItem()
			if item == nil {
				g.Pages.HidePage("modal")
				g.UpdateView()
				return
			}
			switch buttonLabel {
			case "Complete":
				g.td.ChangeStatus(item.id, COMPLETE)
			case "Delete":
				g.td.Delete(item.id)
			}
			g.Pages.HidePage("modal")
			g.UpdateView()
		})
	g.Pages.AddPage("modal", modal, false, false)

	// Listのアイテムを選択した場合の動作を定義
	l.SetSelectedFunc(func(row, column int) {
		item := l.SelectedItem()
		if item != nil {
			g.Pages.ShowPage("modal")
		} else {
			g.UpdateView()
		}
	})
}

// UpdateView 画面描画する
func (l *ListPanel) UpdateView() {
	table := l.Clear()
	for i, item := range l.items {
		cell := tview.NewTableCell(fmt.Sprintf("%4d:%s\n", item.id, item.title))
		table.SetCell(i, 0, cell)
	}
}
