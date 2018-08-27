package main

import (
	"sync"

	oidc "github.com/alextanhongpin/go-openid"
)

// ClientKV represents an in-memory client store.
type ClientKV struct {
	sync.RWMutex
	db map[string]*oidc.Client
}

// NewClientKV returns a pointer to an in-memory client store.
func NewClientKV() *ClientKV {
	return &ClientKV{
		db: make(map[string]*oidc.Client),
	}
}

// Get returns a client by id and a status indicating that the client exist.
func (c *ClientKV) Get(id string) (*oidc.Client, bool) {
	c.RLock()
	client, ok := c.db[id]
	c.RUnlock()
	return client, ok
}

// Put insert a new client by id.
func (c *ClientKV) Put(id string, client *oidc.Client) {
	c.Lock()
	c.db[id] = client
	c.Unlock()
}

// Delete removes a client from the store.
func (c *ClientKV) Delete(id string) {
	c.Lock()
	delete(c.db, id)
	c.Unlock()
}
