package main

import (
	"log"

	"github.com/rivo/tview"
)

// MenuPanel メイン操作パネル
type MenuPanel struct {
	*tview.List
}

// NewMenuPanel creates MenuPanel
func NewMenuPanel() *MenuPanel {
	m := &MenuPanel{
		tview.NewList(),
	}
	m.ShowSecondaryText(false).
		SetBorder(true).               // 枠を表示する
		SetTitle("Menu").              // 枠上のタイトル
		SetTitleAlign(tview.AlignLeft) // タイトルの位置

	return m
}

// Keybinding MenuPanelのキーバインドを設定する.
// Enterを押下した際の動作を定義する.
func (m *MenuPanel) Keybinding(g *GUI) {
	m.AddItem("New Item", "", 'n', func() {
		iform := tview.NewForm()
		iform.AddInputField("Title", "", 20, nil, nil).
			AddInputField("Detail", "", 20, nil, nil).
			AddInputField("Deadline", "", 20, nil, nil).
			AddButton("Save", func() {
				// Saveを選択した場合に、TODOリストにアイテムを追加する
				title := iform.GetFormItem(0).(*tview.InputField).GetText()
				detail := iform.GetFormItem(1).(*tview.InputField).GetText()
				deadline := iform.GetFormItem(2).(*tview.InputField).GetText()
				if err := g.td.Add(title, detail, deadline); err != nil {
					g.App.Stop()
					log.Fatal(err)
				}
				g.UpdateView()
			}).
			AddButton("Cancel", func() {
				g.UpdateView()
			})
		// NewItemを選択した場合に入力フォームを表示する
		g.App.SetRoot(iform, true)
	})
	m.AddItem("Edit", "", 'e', func() {
		g.App.SetFocus(g.ListPanel)
	})
	m.AddItem("Todays", "", 't', func() {
		g.App.SetFocus(g.TodayPanel)
	})
	m.AddItem("Quit", "", 'q', func() {
		// Quitを選択した場合にアプリを終了する
		g.App.Stop()
	})
}
