package todo

import (
	"fmt"
	"time"
	"todocore/entity"
)

// TodoInterface operates TODO
type TodoInterface interface {
	Add(title, detail, deadline string) error
	Delete(id int) error
	ChangeStatus(id, status int) error
	GetActive() ([]Item, error)
	GetComplete() ([]Item, error)
	GetDeadline(deadline int) ([]Item, error)
}

// TODO TODOアプリを管理する
type TODO struct {
	srv      *Server
	ent      entity.Entity
	id       int
	location *time.Location
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
	// DeadlineExpired 期限切れ
	DeadlineExpired = 3

	// SoonSettingStart 「もうすぐ期限」の期間設定
	SoonSettingStart = 3 // 3日前
	// SoonSettingEnd 「もうすぐ期限」の期間設定
	SoonSettingEnd = 1 // 1日前

	// Layout 期日のレイアウト
	Layout = "2006.01.02-15:04"
)

// NewTODO Creates TODO
func NewTODO(ent entity.Entity) *TODO {
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil
	}
	return &TODO{
		ent:      ent,
		location: location,
	}
}

// Add TODOアイテムを追加する
func (td *TODO) Add(title, detail, deadline string) error {
	key := td.ent.NewID()

	// ここでdeadlineをtime.Time形式にコンバートする
	var date time.Time
	if d, err := time.ParseInLocation(Layout, deadline, td.location); err == nil {
		date = d
	}

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

	today := time.Now().Truncate(24 * time.Hour)

	var start, end time.Time
	var entItems []entity.Item
	var err error

	switch deadline {
	// 期限が本日のアイテムを取得する
	case DeadlineToday:
		// 現在の日付を取得して、0:00-23:59を設定する
		start = today
		end = today.Add(time.Hour*23 + time.Minute*59)
		entItems, err = td.ent.GetDate(start, end)

	// 期限が3日以内のアイテムを取得する
	case DeadlineSoon:
		// 3日前の0:00から1日前の23:59を設定する
		start = today.Add(time.Hour * 24 * SoonSettingStart * -1)
		end = today.Add(time.Hour*24*SoonSettingEnd*-1 + time.Minute*59)
		entItems, err = td.ent.GetDate(start, end)

		// 期限切れのアイテムを取得する
	case DeadlineExpired:
		start = today.Add(time.Hour * 24)
		entItems, err = td.ent.GetAfterDate(start)

	default:
		return nil, fmt.Errorf("unsupported definition specified %d", deadline)
	}

	//entItems, err := td.ent.GetDate(start, end)
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
