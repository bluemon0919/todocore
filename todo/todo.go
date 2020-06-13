package todo

import (
	"time"
	"todotool/entity"
)

// TODO TODOアプリを管理する
type TODO struct {
	srv *Server
	ent entity.Entity
	id  int
}

// Item TODOアイテム
type Item struct {
	ID       int       `json:"ID"`
	Title    string    `json:"Title"`
	Detail   string    `json:"Detail"`
	Deadline time.Time `json:"Deadline"`
}

const (
	// ACTIVE TODOアイテムがアクティブ状態
	ACTIVE = 1
	// COMPLETE TODOアイテムが完了状態
	COMPLETE = 2
	// ALL 全てのTODOアイテムを選択するための定義
	ALL = 99

	// DeadlineToday 本日が期限
	DeadlineToday = 1
	// DeadlineSoon もうすぐ期限
	DeadlineSoon = 2

	// SoonSettingStart 「もうすぐ期限」の期間設定
	SoonSettingStart = 3 // 3日前
	// SoonSettingEnd 「もうすぐ期限」の期間設定
	SoonSettingEnd = 1 // 1日前
)

// NewTODO Creates TODO
func NewTODO(ent entity.Entity) *TODO {
	return &TODO{
		ent: ent,
	}
}

// Add TODOアイテムを追加する
func (td *TODO) Add(title, detail string, date time.Time) error {
	key := td.ent.NewID()
	item := &entity.Item{
		Key:    key,
		Title:  title,
		Detail: detail,
		Status: ACTIVE,
		Date:   date,
	}
	td.ent.Add(item)
	return nil
}

// Delete TODOアイテムを削除する
func (td *TODO) Delete(id int) error {
	return td.ent.Delete(id)
}

// ChangeStatus TODOアイテムのステータスを変更する
func (td *TODO) ChangeStatus(id, status int) error {
	return td.ent.Update(id, status)
}

// GetActive Active状態のTODOアイテムを取得する
func (td *TODO) GetActive() ([]Item, error) {
	return td.get(ACTIVE)
}

// GetComplete Complete状態のTODOアイテムを取得する
func (td *TODO) GetComplete() ([]Item, error) {
	return td.get(COMPLETE)
}

// GetDeadline Deadline指定でTODOアイテムを取得する
// startからendの間が期限のアイテムを取得します
func (td *TODO) GetDeadline(deadline int) ([]Item, error) {
	var start, end time.Time
	today := time.Now().Truncate(24 * time.Hour)

	// 期限が本日のアイテムを取得する
	if DeadlineToday == deadline {
		// 現在の日付を取得して、0:00-23:59を設定する
		start = today
		end = today.Add(time.Hour*23 + time.Minute*59)
	}

	// 期限が3日以内のアイテムを取得する
	if DeadlineSoon == deadline {
		// 3日前の0:00から1日前の23:59を設定する
		start = today.Add(time.Hour * 24 * SoonSettingStart * -1)
		end = today.Add(time.Hour*24*SoonSettingEnd*-1 + time.Minute*59)
	}
	entItems, err := td.ent.GetDate(start, end)
	if err != nil {
		return nil, err
	}
	var items []Item
	for _, eitem := range entItems {
		items = append(items, Item{
			ID:       eitem.Key,
			Title:    eitem.Title,
			Detail:   eitem.Detail,
			Deadline: eitem.Date,
		})
	}
	return items, nil
}

func (td *TODO) get(status int) ([]Item, error) {
	var items []Item
	eis, err := td.ent.Get(status)
	if err != nil {
		return nil, err
	}
	for _, ei := range eis {
		items = append(items, Item{
			ID:       ei.Key,
			Title:    ei.Title,
			Detail:   ei.Detail,
			Deadline: ei.Date,
		})
	}
	return items, nil
}
