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
func (c *ClientKV) Get(name string) (*oidc.Client, bool) {
	c.RLock()
	client, ok := c.db[name]
	c.RUnlock()
	return client, ok
}

// GetByID returns the client by client id.
func (c *ClientKV) GetByID(id string) (client *oidc.Client) {
	c.RLock()
	for _, c := range c.db {
		if c.ClientID == id {
			client = c
			break
		}
	}
	c.RUnlock()
	return
}

// GetByIDAndSecret returns the client by client id and client secret.
func (c *ClientKV) GetByIDAndSecret(id, secret string) (client *oidc.Client) {
	c.RLock()
	for _, c := range c.db {
		if c.ClientID == id && c.ClientSecret == secret {
			client = c
			break
		}
	}
	c.RUnlock()
	return
}

// Put insert a new client by id.
func (c *ClientKV) Put(name string, client *oidc.Client) {
	c.Lock()
	c.db[name] = client
	c.Unlock()
}

// Delete removes a client from the store.
func (c *ClientKV) Delete(name string) {
	c.Lock()
	delete(c.db, name)
	c.Unlock()
}
