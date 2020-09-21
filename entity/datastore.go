package entity

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
)

// EntityDatastore データを管理する
type EntityDatastore struct {
	id         int
	projectID  string
	entityType string
}

// NewDatastore creates Entity
func NewDatastore(projectID string, entityType string) *EntityDatastore {
	return &EntityDatastore{
		projectID:  projectID,
		entityType: entityType,
	}
}

// Add Entityにアイテムを追加する
func (ent *EntityDatastore) Add(item *Item) error {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, ent.projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	dkey := datastore.NameKey(ent.entityType, "", nil)
	_, err = client.Put(ctx, dkey, item)
	return err
}

// Delete Entityから指定キーを削除する
func (ent *EntityDatastore) Delete(id int64) error {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, ent.projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	return client.Delete(ctx, &datastore.Key{
		Kind: ent.entityType,
		ID:   id,
	})
}

// Update Entityの指定のキーを入力ステータスでアップデートする
func (ent *EntityDatastore) Update(id int64, status int) error {

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, ent.projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	var item Item
	key := &datastore.Key{
		Kind: ent.entityType,
		ID:   id,
	}
	if err := client.Get(ctx, key, &item); err != nil {
		return err
	}
	item.Status = status
	if _, err := client.Put(ctx, key, &item); err != nil {
		return err
	}
	return nil
}

// Get Entityからアイテムを取得する
func (ent *EntityDatastore) Get(status int) (items []Item, ids []int64, err error) {
	query := datastore.NewQuery(ent.entityType).Filter("Status = ", status)
	items, ids, err = ent.get(query)
	return
}

// GetDate 期間を指定してアイテムを取得する
func (ent *EntityDatastore) GetDate(start, end time.Time) (items []Item, ids []int64, err error) {
	query := datastore.NewQuery(ent.entityType).Filter("EntTime >= ", start).Filter("EntTime <= ", end)
	items, ids, err = ent.get(query)
	return
}

// GetAfterDate 期間を指定してアイテムを取得する
func (ent *EntityDatastore) GetAfterDate(base time.Time) (items []Item, ids []int64, err error) {
	query := datastore.NewQuery(ent.entityType)
	tmpItems, tmpIds, err := ent.get(query)
	if err != nil {
		return nil, nil, err
	}
	for i, item := range tmpItems {
		// base <= item.EndTime
		if !base.After(item.EndTime) {
			items = append(items, item)
			ids = append(ids, tmpIds[i])
		}
	}
	return
}

// GetBeforeDate は基準日時以前のアイテムを取得する
func (ent *EntityDatastore) GetBeforeDate(base time.Time) (items []Item, ids []int64, err error) {
	query := datastore.NewQuery(ent.entityType)
	tmpItems, tmpIds, err := ent.get(query)
	if err != nil {
		return nil, nil, err
	}
	for i, item := range tmpItems {
		// base >= item.EndTime
		if !base.Before(item.EndTime) {
			items = append(items, item)
			ids = append(ids, tmpIds[i])
		}
	}
	return
}

// get Entityからアイテムを取得する
func (ent *EntityDatastore) get(query *datastore.Query) (items []Item, ids []int64, err error) {

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, ent.projectID)
	if err != nil {
		return nil, nil, err
	}
	defer client.Close()

	var keys []*datastore.Key
	if keys, err = client.GetAll(ctx, query, &items); err != nil {
		return nil, nil, err
	}
	for _, key := range keys {
		ids = append(ids, key.ID)
	}
	return
}
