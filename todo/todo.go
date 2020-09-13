package todo

import (
	"fmt"
	"time"

	"github.com/bluemon0919/todocore/entity"
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
	ID        int64
	Title     string    `json:"Title"`
	Detail    string    `json:"Detail"`
	StartTime time.Time `json:"StartTime"`
	Deadline  time.Time `json:"Deadline"`
	StationID string    `json:"StationID"`
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
func (td *TODO) Add(title, detail, startTime, deadline, stationID string) error {
	var end time.Time
	if d, err := time.ParseInLocation(Layout, deadline, td.location); err == nil {
		end = d
	}
	var start time.Time
	if d, err := time.ParseInLocation(Layout, startTime, td.location); err == nil {
		start = d
	}

	item := &entity.Item{
		Title:     title,
		Detail:    detail,
		Status:    ACTIVE,
		StartTime: start,
		EndTime:   end,
		StationID: stationID,
	}
	td.ent.Add(item)
	return nil
}

// Delete TODOアイテムを削除する
func (td *TODO) Delete(id int64) error {
	return td.ent.Delete(id)
}

// ChangeStatus TODOアイテムのステータスを変更する
func (td *TODO) ChangeStatus(id int64, status int) error {
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

// GetDeadline はDeadline指定でTODOアイテムを取得する
// startからendの間が期限のアイテムを取得します
func (td *TODO) GetDeadline(deadline int) ([]Item, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}
	now := time.Now().In(loc)
	today := now.Truncate(time.Hour).Add(-time.Duration(now.Hour()) * time.Hour)

	var entItems []entity.Item
	var ids []int64

	switch deadline {
	// 期限が本日のアイテムを取得する
	case DeadlineToday:
		// 現在の日付を取得して、0:00-23:59を設定する
		start := today
		end := today.Add(time.Hour*23 + time.Minute*59)
		entItems, ids, err = td.ent.GetDate(start, end)

	// 期限が3日以内のアイテムを取得する
	case DeadlineSoon:
		// 3日前の0:00から1日前の23:59を設定する
		start := today.Add(time.Hour * 24 * SoonSettingStart * -1)
		end := today.Add(time.Hour*24*SoonSettingEnd*-1 + time.Minute*59)
		entItems, ids, err = td.ent.GetDate(start, end)

		// 期限切れのアイテムを取得する
	case DeadlineExpired:
		base := today
		entItems, ids, err = td.ent.GetBeforeDate(base)

	default:
		return nil, fmt.Errorf("unsupported definition specified %d", deadline)
	}

	if err != nil {
		return nil, err
	}
	if len(entItems) != len(ids) {
		return nil, fmt.Errorf("the number of items and the number of keys are different")
	}
	var items []Item
	for i, eitem := range entItems {
		items = append(items, Item{
			ID:        ids[i],
			Title:     eitem.Title,
			Detail:    eitem.Detail,
			StartTime: eitem.StartTime,
			Deadline:  eitem.EndTime,
			StationID: eitem.StationID,
		})
	}
	return items, nil
}

func (td *TODO) get(status int) ([]Item, error) {
	var items []Item
	eis, ids, err := td.ent.Get(status)
	if err != nil {
		return nil, err
	}
	if len(eis) != len(ids) {
		return nil, fmt.Errorf("the number of items and the number of keys are different")
	}
	for i, ei := range eis {
		items = append(items, Item{
			ID:        ids[i],
			Title:     ei.Title,
			Detail:    ei.Detail,
			StartTime: ei.StartTime.In(td.location),
			Deadline:  ei.EndTime.In(td.location),
			StationID: ei.StationID,
		})
	}
	return items, nil
}
