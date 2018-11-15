package main

import "time"

type Code struct {
	CreatedAt time.Time
	ID        string
	TTL       time.Duration
}

func NewCode(id string, ttl time.Duration) *Code {
	return &Code{
		ID:        id,
		TTL:       ttl,
		CreatedAt: time.Now(),
	}
}

func (c *Code) HasExpired() bool {
	return time.Since(c.CreatedAt) > c.TTL
}

func GetCode(repo CodeRepository, id string) (*Code, error) {
	return repo.GetCodeByID(id)
}
