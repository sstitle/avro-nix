package item

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RpcClient is an adapter that implements Repository over HTTP+JSON.
// It delegates all calls to a remote server, making local↔remote swap
// transparent to the application core and service layer.
//
// The wire protocol mirrors the messages in ports/item/repository.avpr.
type RpcClient struct {
	base   string
	client *http.Client
}

func NewRpcClient(addr string) *RpcClient {
	return &RpcClient{base: "http://" + addr, client: &http.Client{}}
}

func (c *RpcClient) post(path string, body, out any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	resp, err := c.client.Post(c.base+path, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		var nf NotFoundError
		json.NewDecoder(resp.Body).Decode(&nf)
		return &nf
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server error: %s", resp.Status)
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

func (c *RpcClient) Get(id string) (*Item, error) {
	var result Item
	err := c.post("/item/get", map[string]string{"id": id}, &result)
	return &result, err
}

func (c *RpcClient) List() ([]*Item, error) {
	var result []*Item
	err := c.post("/item/list", struct{}{}, &result)
	return result, err
}

func (c *RpcClient) Save(it *Item) error {
	return c.post("/item/save", it, nil)
}

func (c *RpcClient) Delete(id string) error {
	return c.post("/item/delete", map[string]string{"id": id}, nil)
}
