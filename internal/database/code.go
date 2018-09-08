package database

import (
	"sync"

	oidc "github.com/alextanhongpin/go-openid"
)

type CodeKV struct {
	sync.RWMutex
	db map[string]*oidc.Code
}

func NewCodeKV() *CodeKV {
	return &CodeKV{
		db: make(map[string]*oidc.Code),
	}
}

func (c *CodeKV) Get(id string) (*oidc.Code, bool) {
	c.RLock()
	code, ok := c.db[id]
	c.RUnlock()
	return code, ok
}

func (c *CodeKV) Put(id string, code *oidc.Code) {
	c.Lock()
	c.db[id] = code
	c.Unlock()
}

func (c *CodeKV) Delete(id string) {
	c.Lock()
	delete(c.db, id)
	c.Unlock()
}
