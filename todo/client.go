package todo

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Client struct {
	url string
}

// NewClient creates new client.
// Gives the URL of the server.
func NewClient(url string) (*Client, error) {
	if len(url) == 0 {
		return nil, errors.New("missing url")
	}
	return &Client{
		url: url,
	}, nil
}

// Add TODOアイテムを追加する
func (c *Client) Add(title, detail string) error {
	datas := []AddRequest{
		{
			Title:  title,
			Detail: detail,
		},
	}
	bs, err := json.Marshal(datas)
	if err != nil {
		return err
	}

	body := bytes.NewReader(bs)
	url := c.url + "/?kind=add"
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return err
	}

	client := new(http.Client)
	_, err = client.Do(req)
	return err
}

// Delete TODOアイテムを削除する
func (c *Client) Delete(id int) error {
	ids := []int{id}
	bs, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	body := bytes.NewReader(bs)
	req, err := http.NewRequest(http.MethodDelete, c.url, body)
	if err != nil {
		return err
	}

	client := new(http.Client)
	_, err = client.Do(req)
	return err
}

// ChangeStatus TODOアイテムのステータスを変更する
func (c *Client) ChangeStatus(id, status int) error {
	data := []UpdateRequest{
		{id, status},
	}
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}

	body := bytes.NewReader(bs)
	req, err := http.NewRequest(http.MethodPut, c.url+"/?kind=update", body)
	if err != nil {
		return err
	}

	client := new(http.Client)
	_, err = client.Do(req)
	return err
}

// GetActive Active状態のTODOアイテムを取得する
func (c *Client) GetActive() ([]Item, error) {
	url := c.url + "/?kind=active"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res Response
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&res); err != nil {
		return nil, err
	}
	return res.Items, nil
}

// GetComplete Complete状態のTODOアイテムを取得する
func (c *Client) GetComplete() ([]Item, error) {
	url := c.url + "/?kind=complete"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res Response
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&res); err != nil {
		return nil, err
	}
	return res.Items, nil
}
