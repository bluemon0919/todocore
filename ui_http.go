package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"todocore/todo"
)

// Handler HTTPハンドラおよびHTML形式のユーザーインターフェースを提供する
type Handler struct {
	addr string
	td   *todo.Client // TODOツールへの操作
}

// NewHandler create WebHandler
func NewHandler(td *todo.Client, addr string) *Handler {
	return &Handler{
		addr: addr,
		td:   td,
	}
}

// Run 実行する
func (h *Handler) Run() error {
	// httpハンドラの設定
	http.HandleFunc("/list", h.listHandler)
	http.HandleFunc("/operation_post", h.operationHandler)
	http.HandleFunc("/input", h.inputFormHandler)
	http.HandleFunc("/input_post", h.inputPostHandler)

	// httpサーバ起動
	return http.ListenAndServe(h.addr, nil)
}

// listHandler
// TODOアイテム一覧の表示と入力フォームへの遷移を行う
func (h *Handler) listHandler(w http.ResponseWriter, r *http.Request) {

	items, _ := h.td.GetActive()

	type ListItem struct {
		Key   int
		Title string
	}

	var listItems []ListItem
	for _, item := range items {
		listItem := ListItem{
			Key:   item.ID,
			Title: item.Title,
		}
		listItems = append(listItems, listItem)
	}

	tpl := template.Must(template.ParseFiles("static/list.html"))
	tpl.Execute(w, listItems)
}

// operationHandler
// TODOアイテム一覧からの操作を処理する
func (h *Handler) operationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/list", http.StatusFound)
		return
	}

	// 入力テキストをレコードに登録する
	keyStr := r.FormValue("select_key")
	key, err := strconv.Atoi(keyStr)
	if err != nil {
		http.Redirect(w, r, "/list", http.StatusFound)
		return
	}

	if "" != r.FormValue("complete") {
		h.td.ChangeStatus(key, todo.COMPLETE)
	}
	if "" != r.FormValue("delete") {
		h.td.Delete(key)
	}
	http.Redirect(w, r, "/list", http.StatusFound)
}

// inputFormHandler
// 入力フォームを表示する
func (h *Handler) inputFormHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("static/input.html"))
	tpl.Execute(w, nil)
}

// inputPostHandler
// 入力フォームからのPOST処理を行う
func (h *Handler) inputPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/list", http.StatusFound)
		return
	}
	// 入力テキストをレコードに登録する
	title := r.FormValue("title")
	detail := r.FormValue("detail")
	deadline := r.FormValue("deadline")

	h.td.Add(title, detail, deadline)
	http.Redirect(w, r, "/list", http.StatusFound)
}
