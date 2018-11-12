package client

import openid "github.com/alextanhongpin/go-openid"

// Nop represents a no-op service for testing edge cases or mocking the service.
type Nop struct{}

func (n *Nop) Register(c *openid.Client) (*openid.Client, error) {
	return c, nil
}

func (n *Nop) Read(clientID string) (*openid.Client, error) {
	return nil, nil
}
