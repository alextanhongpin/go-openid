package openid

import "sync"

// ClientStore implements an in-memory ClientStore for testing
type ClientStore struct {
	sync.RWMutex
	db map[string]*Client
}

func NewClientStore() *ClientStore {
	return &ClientStore{
		db: make(map[string]*Client),
	}
}

func (c *ClientStore) Get(id string) *Client {
	c.RLock()
	client, ok := c.db[id]
	c.RUnlock()
	if !ok {
		return nil
	}
	return client
}

func (c *ClientStore) Put(id string, client *Client) {
	c.Lock()
	c.db[id] = client
	c.Unlock()
}

func (c *ClientStore) Delete(id string) {
	c.Lock()
	delete(c.db, id)
	c.Unlock()
}
