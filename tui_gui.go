package main

import (
	"todotool/todo"

	"github.com/rivo/tview"
)

// GUI ユーザーインターフェース
type GUI struct {
	App        *tview.Application // tview全体を制御する. Run()でTUIツールを起動する
	Pages      *tview.Pages       // 各画面へのフォーカスを行う
	MenuPanel  *MenuPanel         // メニュー画面
	ListPanel  *ListPanel
	TodayPanel *ListPanel
	td         *todo.Client // TODOツールへの操作
}

// NewGUI creates GUI
func NewGUI(td *todo.Client) *GUI {
	return &GUI{
		App:        tview.NewApplication(),
		Pages:      tview.NewPages(),
		MenuPanel:  NewMenuPanel(),
		ListPanel:  NewListPanel("ToDo List"),
		TodayPanel: NewListPanel("Today limit"),
		td:         td,
	}
}

// Run 実行する
func (g *GUI) Run() error {

	g.UpdateView()

	grid := tview.NewGrid(). // グリッドを作成する
					SetColumns(20, 0).                            // グリッドのセルを縦に二つ作る
					AddItem(g.MenuPanel, 0, 0, 2, 1, 0, 0, true). // MenuPanelを左側に配置
					AddItem(g.ListPanel, 0, 1, 1, 1, 0, 0, true). // ListPanelを右側に配置
					AddItem(g.TodayPanel, 1, 1, 1, 1, 0, 0, true) // TodayPanelを右側に配置

	g.Pages.AddAndSwitchToPage("main", grid, true) // グリッドをPagesの管理下に置き、フォーカスする
	g.SetKeybinding()                              // 各画面のキーバインドを設定する

	return g.App.SetRoot(g.Pages, true).Run() // PagesをApplicationの管理下におき、Run()でtviewを実行する
}

// SetKeybinding キーバインドを設定する
func (g *GUI) SetKeybinding() {
	g.MenuPanel.Keybinding(g)
	g.ListPanel.Keybinding(g)
	g.TodayPanel.Keybinding(g)
}

// UpdateView メイン画面を更新する
func (g *GUI) UpdateView() {
	items, _ := g.td.GetActive()
	g.ListPanel.SetItems(items)
	g.ListPanel.UpdateView()

	todays, _ := g.td.GetDeadlineToday()
	g.TodayPanel.SetItems(todays)
	g.TodayPanel.UpdateView()

	g.App.SetRoot(g.Pages, true)
}
