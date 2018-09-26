package database

import (
	"errors"
	"sync"

	oidc "github.com/alextanhongpin/go-openid"
)

// ClientKV represents an in-memory client store.
type ClientKV struct {
	sync.RWMutex
	db map[string]*oidc.Client
	// TODO: Add cache-layer, both in-memory and redis
}

// NewClientKV returns a pointer to an in-memory client store.
func NewClientKV() *ClientKV {
	return &ClientKV{
		db: make(map[string]*oidc.Client),
	}
}

// Get returns a client by id and a status indicating that the client exist.
func (c *ClientKV) Get(id string) (*oidc.Client, error) {
	c.RLock()
	client, ok := c.db[id]
	c.RUnlock()
	if !ok {
		return nil, errors.New("client does not exist")
	}
	return client, nil
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

// Has returns true if the client id exist in the storage.
func (c *ClientKV) Has(id string) bool {
	c.RLock()
	_, ok := c.db[id]
	c.RUnlock()
	return ok
}

// GetByCredentials returns the client by client id and client secret.
func (c *ClientKV) GetByCredentials(clientID, clientSecret string) (client *oidc.Client, err error) {
	c.RLock()
	for _, c := range c.db {
		if c.ClientID == clientID && c.ClientSecret == clientSecret {
			client = c
			break
		}
	}
	c.RUnlock()
	return
}

// Put insert a new client by id.
func (c *ClientKV) Put(id string, client *oidc.Client) error {
	c.Lock()
	c.db[id] = client
	c.Unlock()
	return nil
}

// Delete removes a client from the store.
func (c *ClientKV) Delete(id string) {
	c.Lock()
	delete(c.db, id)
	c.Unlock()
}

// List returns a paginated list of users.
func (c *ClientKV) List(limit int) []*oidc.Client {
	var i int
	var clients []*oidc.Client
	for _, v := range c.db {
		if i == limit {
			break
		}
		clients = append(clients, v)
		i++
	}
	return clients
}
