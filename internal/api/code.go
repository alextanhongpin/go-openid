package main

import "time"

type (
	Code struct {
		CreatedAt time.Time
		ID        string
		TTL       time.Duration
	}
	CodeBuilder struct {
		defaults Code
		// This is more useful than writable, if there are values
		// that depends on the code produced.
		overwrite func(c *Code)
		// writable bool
	}
)

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

func CreateCode(repo CodeRepository, code *Code) error {
	return repo.Create(code)
}

type CodeFactory func() *Code
