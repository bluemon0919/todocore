package entity

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
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

// NewID 新しいIDを返す
func (ent *EntityDatastore) NewID() int {
	ent.id++
	return ent.id
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
func (ent *EntityDatastore) Delete(key int) error {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, ent.projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	query := datastore.NewQuery(ent.entityType).Filter("Key = ", key)
	it := client.Run(ctx, query)
	var item Item
	dkey, err := it.Next(&item)
	if err == iterator.Done {
		return err
	}

	return client.Delete(ctx, dkey)
}

// Update Entityの指定のキーを入力ステータスでアップデートする
func (ent *EntityDatastore) Update(key, status int) error {

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, ent.projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	query := datastore.NewQuery(ent.entityType).Filter("Key = ", key)
	it := client.Run(ctx, query)

	var item Item
	dkey, err := it.Next(&item)
	if err == iterator.Done {
		return err
	}
	item.Status = status
	_, err = client.Put(ctx, dkey, &item)
	return err
}

// Get Entityからアイテムを取得する
func (ent *EntityDatastore) Get(status int) (items []Item, err error) {
	query := datastore.NewQuery(ent.entityType).Filter("Status = ", status)
	items, err = ent.get(query)
	return
}

// GetDate 期間を指定してアイテムを取得する
func (ent *EntityDatastore) GetDate(start, end time.Time) (items []Item, err error) {
	query := datastore.NewQuery(ent.entityType).Filter("EntTime >= ", start).Filter("EntTime <= ", end)
	items, err = ent.get(query)
	return
}

// GetAfterDate 期間を指定してアイテムを取得する
func (ent *EntityDatastore) GetAfterDate(base time.Time) (items []Item, err error) {
	query := datastore.NewQuery(ent.entityType).Filter("EntTime >= ", base)
	items, err = ent.get(query)
	return
}

// GetBeforeDate は基準日時以前のアイテムを取得する
func (ent *EntityDatastore) GetBeforeDate(base time.Time) (items []Item, err error) {
	query := datastore.NewQuery(ent.entityType).Filter("EntTime <= ", base)
	items, err = ent.get(query)
	return
}

// get Entityからアイテムを取得する
func (ent *EntityDatastore) get(query *datastore.Query) (items []Item, err error) {

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, ent.projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	if _, err := client.GetAll(ctx, query, &items); err != nil {
		return nil, err
	}

	// Keyでソートする
	keysort := func(p1, p2 *Item) bool {
		return p1.Key < p2.Key
	}
	By(keysort).Sort(items)
	return
}
