package main

import (
	"github.com/rivo/tview"
)

// GUI ユーザーインターフェース
type GUI struct {
	App         *tview.Application // tview全体を制御する. Run()でTUIツールを起動する
	Pages       *tview.Pages       // 各画面へのフォーカスを行う
	MenuPanel   *MenuPanel         // メニュー画面
	ListPanel   *ListPanel
	DetailPanel *DetailPanel
	td          *TODO // TODOツールへの操作
}

// NewGUI creates GUI
func NewGUI(td *TODO) *GUI {
	return &GUI{
		App:         tview.NewApplication(),
		Pages:       tview.NewPages(),
		MenuPanel:   NewMenuPanel(),
		ListPanel:   NewListPanel(),
		DetailPanel: NewDetailPanel(),
		td:          td,
	}
}

// Run 実行する
func (g *GUI) Run() error {

	g.UpdateView()

	grid := tview.NewGrid(). // グリッドを作成する
					SetColumns(20, 0).                            // グリッドのセルを縦に二つ作る
					AddItem(g.MenuPanel, 0, 0, 1, 1, 0, 0, true). // MenuPanelを左側に配置
					AddItem(g.ListPanel, 0, 1, 1, 1, 0, 0, true)  // ListPanelを右側に配置

	g.Pages.AddAndSwitchToPage("main", grid, true) // グリッドをPagesの管理下に置き、フォーカスする
	g.SetKeybinding()                              // 各画面のキーバインドを設定する

	return g.App.SetRoot(g.Pages, true).Run() // PagesをApplicationの管理下におき、Run()でtviewを実行する
}

// SetKeybinding キーバインドを設定する
func (g *GUI) SetKeybinding() {
	g.MenuPanel.Keybinding(g)
	g.ListPanel.Keybinding(g)
}

// UpdateView メイン画面を更新する
func (g *GUI) UpdateView() {
	items, _ := g.td.GetActive()
	g.ListPanel.SetItems(items)
	g.ListPanel.UpdateView()

	g.App.SetRoot(g.Pages, true)
}
