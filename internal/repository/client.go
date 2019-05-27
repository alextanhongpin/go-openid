package repository

import (
	"sync"

	openid "github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/randstr"
)

type Client struct {
	sync.RWMutex
	clients map[string]*openid.Client
}

func NewClient() *Client {
	return &Client{
		clients: make([]*openid.Client),
	}
}

// WithCredentials returns the client by client id and client secret.
func (c *ClientKV) WithCredentials(clientID, clientSecret string) (client *openid.Client, err error) {
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

// Create insert a new client by id.
func (c *ClientKV) Create(client openid.Client) error {
	id := randstr.RandomString()
	c.Lock()
	// Generate a random id.
	c.db[id] = &client
	c.Unlock()
	return id, nil
}
