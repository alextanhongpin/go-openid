package main

import (
	"sync"

	openid "github.com/alextanhongpin/go-openid"
)

type ClientKV struct {
	sync.RWMutex
	db map[string]*openid.Client
}

func NewClientKV() *ClientKV {
	return &ClientKV{
		db: make(map[string]*openid.Client),
	}
}

func (c *ClientKV) Get(id string) (*openid.Client, bool) {
	c.RLock()
	client, ok := c.db[id]
	c.RUnlock()
	return client, ok
}

func (c *ClientKV) Put(id string, client *openid.Client) {
	c.Lock()
	c.db[id] = client
	c.Unlock()
}

func (c *ClientKV) Delete(id string) {
	c.Lock()
	delete(c.db, id)
	c.Unlock()
}
