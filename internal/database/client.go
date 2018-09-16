package database

import (
	"sync"

	oidc "github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
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

func (c *ClientKV) Has(id string) bool {
	c.RLock()
	_, ok := c.db[id]
	c.RUnlock()
	return ok
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

func (c *ClientKV) GenerateClientCredentials() (clientID, clientSecret string) {
	// for i := 0; i < defaultRetry; i++ {
	//         clientID = crypto.NewXID()
	//         if exist := c.Has(clientID); !exist {
	//                 break
	//         }
	// }
	clientID = crypto.NewXID()
	var err error
	clientSecret, err = crypto.GenerateRandomString(32)
	if err != nil {
		panic(err)
	}
	return
}
