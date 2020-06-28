package main

import (
	"fmt"
	"math/rand"
	"todocore/todo"

	"github.com/rivo/tview"
)

// ListPanel ファイル操作するための情報
type ListPanel struct {
	items []todo.Item
	*tview.Table
	modelName string
	title     string
}

// NewListPanel creates ListPanel
func NewListPanel(title string) *ListPanel {
	p := &ListPanel{
		Table:     tview.NewTable(),
		modelName: RandString(10),
		title:     title,
	}

	p.SetBorder(true).
		SetTitle(title).
		SetTitleAlign(tview.AlignLeft)

	p.SetSelectable(true, false)

	return p
}

// RandString ランダムな文字列を作成する
func RandString(n int) string {
	const strset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = strset[rand.Intn(len(strset))]
	}
	return string(b)
}

// SetItems set files info
func (l *ListPanel) SetItems(items []todo.Item) {
	l.items = items
}

// SelectedItem 現在選択しているファイル情報を取得する
func (l *ListPanel) SelectedItem() *todo.Item {
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
				g.Pages.HidePage(l.modelName)
				g.UpdateView()
				return
			}
			switch buttonLabel {
			case "Complete":
				g.td.ChangeStatus(item.ID, todo.COMPLETE)
			case "Delete":
				g.td.Delete(item.ID)
			}
			g.Pages.HidePage(l.modelName)
			g.UpdateView()
		})
	g.Pages.AddPage(l.modelName, modal, false, false)

	// Listのアイテムを選択した場合の動作を定義
	l.SetSelectedFunc(func(row, column int) {
		item := l.SelectedItem()
		if item != nil {
			g.Pages.ShowPage(l.modelName)
		} else {
			g.UpdateView()
		}
	})
}

// UpdateView 画面描画する
func (l *ListPanel) UpdateView() {
	table := l.Clear()
	for i, item := range l.items {
		sTitle := item.Title[:min(10, len(item.Title))]
		deadline := item.Deadline.Format(todo.Layout)
		sDeadline := deadline[:min(16, len(deadline))]
		cell := tview.NewTableCell(
			fmt.Sprintf("%4d:%-20s %-16s\n",
				item.ID, sTitle, sDeadline))
		table.SetCell(i, 0, cell)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
