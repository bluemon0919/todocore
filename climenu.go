package main

import (
	"fmt"
	"io"
	"todocore/todo"
)

// Menu TODOリストを操作するメニュー
type Menu struct {
	td    todo.TodoInterface
	stdin io.Reader
}

// NewMenu コマンドラインメニューでTODOを操作する
func NewMenu(td todo.TodoInterface, stdin io.Reader) *Menu {
	return &Menu{
		td:    td,
		stdin: stdin,
	}
}

// Run コマンドラインメニューを実行します
func (m *Menu) Run() error {
	for {
		ShowMenu()
		selectedMenu := 0
		if _, err := fmt.Fscanf(m.stdin, "%d", &selectedMenu); err != nil {
			return err
		}

		var err error
		switch selectedMenu {
		case 1:
			err = m.NewIssue()
		case 2:
			err = m.SelectIssue()
		case 3:
			err = m.ShowComplete()
		case 9:
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// ShowMenu メインメニューを表示する
func ShowMenu() {
	fmt.Println("=== Menu ===")
	fmt.Println("= 1:New    =")
	fmt.Println("= 2:List   =")
	fmt.Println("= 3:CmpList=")
	fmt.Println("= 9:End    =")
	fmt.Println("============")
}

// ShowSelectMenu Selectのメニューを表示する
func ShowSelectMenu() {
	fmt.Println("=== Select Menu ===")
	fmt.Println("= 1:Complete      =")
	fmt.Println("= 2:Delete        =")
	fmt.Println("= 9:End           =")
	fmt.Println("===================")
}

// NewIssue TODOへの追加を行う
func (m *Menu) NewIssue() error {
	fmt.Println("Please input a new todo issue.")

	fmt.Print("title:")
	title := ""
	if _, err := fmt.Fscanf(m.stdin, "%s", &title); err != nil {
		return err
	}

	fmt.Print("detail:")
	detail := ""
	if _, err := fmt.Fscanf(m.stdin, "%s", &detail); err != nil {
		return err
	}

	fmt.Print("deadline:")
	deadline := ""
	if _, err := fmt.Fscanf(m.stdin, "%s", &deadline); err != nil {
		return err
	}
	return m.td.Add(title, detail, deadline)
}

// SelectIssue TODOのリストを表示して処理を選択させる
func (m *Menu) SelectIssue() error {

	// 登録しているTODOアイテムを表示する
	fmt.Println("TODO List.")
	items, err := m.td.GetActive()
	if err != nil {
		return err
	}
	for _, item := range items {
		fmt.Printf("%d: %s\n", item.ID, item.Title)
	}

	// 表示したリストから、アイテムを選択する
	fmt.Println("Please select a item. Select 0 to return.")
	var n int
	if _, err := fmt.Fscanf(m.stdin, "%d", &n); err != nil {
		return err
	}
	if n == 0 {
		return nil
	}

	// 選択したアイテムの操作を選択する
	fmt.Println("Please select an operation")
	ShowSelectMenu()
	var op int
	if _, err := fmt.Fscanf(m.stdin, "%d", &op); err != nil {
		return err
	}
	switch op {
	case 1:
		err = m.td.ChangeStatus(n, todo.COMPLETE)
	case 2:
		err = m.td.Delete(n)
	}
	return err
}

// ShowComplete 完了しているアイテムを表示する
func (m *Menu) ShowComplete() error {
	fmt.Println("Show Complete List.")
	items, err := m.td.GetComplete()
	if err != nil {
		return err
	}

	for _, item := range items {
		fmt.Printf("%d: %s\n", item.ID, item.Title)
	}
	return nil
}
