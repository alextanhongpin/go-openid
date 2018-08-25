package openid

import (
	"sync"
	"time"
)

// Code represents the authorization code
type Code struct {
	Code      string
	CreatedAt time.Time
	TTL       time.Duration
}

// Expired returns if the code has reached pass the expiration limit
func (c *Code) Expired() bool {
	return time.Since(c.CreatedAt) > c.TTL
}

// CodeStore keeps a cache of the authorization code locally
// Replace it with a distributed cache when running multiple nodes
type CodeStore struct {
	sync.RWMutex
	codes map[string]*Code
}

// NewCodeStore returns a pointer to a new CodeStore
func NewCodeStore() *CodeStore {
	return &CodeStore{
		codes: make(map[string]*Code),
	}
}

// Get returns a pointer to the Code from the cache from the given key
func (s *CodeStore) Get(id string) *Code {
	s.RLock()
	code, ok := s.codes[id]
	s.RUnlock()
	if !ok {
		return nil
	}
	return code
}

// Put stores the code in the cache for the given key
func (s *CodeStore) Put(id string, code *Code) {
	s.Lock()
	s.codes[id] = code
	s.Unlock()
}

// Delete removes the code in the cache for the given key
func (s *CodeStore) Delete(id string) {
	s.Lock()
	delete(s.codes, id)
	s.Unlock()
}
