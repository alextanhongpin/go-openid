package main

import (
	"sync"

	openid "github.com/alextanhongpin/go-openid"
)

type ClientKV struct {
	sync.RWMutex
	db map[string]*openid.Client
}

func (c *ClientKV) Get(id string) *openid.Client {
	if id == "" {
		return nil
	}
	c.RLock()
	client, ok := c.db[id]
	c.RUnlock()
	if !ok {
		return nil
	}
	return client
}

func (c *ClientKV) Put(id string, client *openid.Client) {
	if id == "" || client == nil {
		return
	}
	c.Lock()
	c.db[id] = client
	c.Unlock()
}

func (c *ClientKV) Delete(id string) {
	if id == "" {
		return
	}
	c.Lock()
	delete(c.db, id)
	c.Unlock()
}
