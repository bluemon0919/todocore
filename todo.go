package main

// TODO TODOアプリを管理する
type TODO struct {
	e  *Entity
	id int
}

// Item TODOアイテム
type Item struct {
	id     int
	title  string
	detail string
}

const (
	// ACTIVE TODOアイテムがアクティブ状態
	ACTIVE = 1
	// COMPLETE TODOアイテムが完了状態
	COMPLETE = 2
	// ALL 全てのTODOアイテムを選択するための定義
	ALL = 99
)

// NewTODO Creates TODO
func NewTODO(e *Entity) *TODO {
	return &TODO{
		e: e,
	}
}

// Add TODOアイテムを追加する
func (td *TODO) Add(title, detail string) error {
	ei := &EntityItem{
		key:    td.ID(),
		title:  title,
		detail: detail,
		status: ACTIVE,
	}
	return td.e.Add(ei)
}

// Delete TODOアイテムを削除する
func (td *TODO) Delete(id int) error {
	return td.e.Delete(id)
}

// ChangeStatus TODOアイテムのステータスを変更する
func (td *TODO) ChangeStatus(id, status int) error {
	return td.e.Update(id, status)
}

// GetActive Active状態のTODOアイテムを取得する
func (td *TODO) GetActive() ([]Item, error) {
	return td.get(ACTIVE)
}

// GetComplete Complete状態のTODOアイテムを取得する
func (td *TODO) GetComplete() ([]Item, error) {
	return td.get(COMPLETE)
}

func (td *TODO) get(kind int) ([]Item, error) {
	var items []Item
	eis, err := td.e.Get(kind)
	if err != nil {
		return nil, err
	}
	for _, ei := range eis {
		items = append(items, Item{
			ei.key,
			ei.title,
			ei.detail,
		})
	}
	return items, nil
}

// ID 固有のIDを発行する
func (td *TODO) ID() int {
	td.id++
	return td.id
}
