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

func NewCodeBuilder(defaults Code) *CodeBuilder {
	return &CodeBuilder{
		defaults:  defaults,
		overwrite: func(c *Code) {},
	}
}

func (c *CodeBuilder) Setoverwrite(overwrite func(c *Code)) {
	c.overwrite = overwrite
}

func (c *CodeBuilder) SetCreatedAt(createdAt time.Time) {
	c.defaults.CreatedAt = createdAt
}

func (c *CodeBuilder) SetID(id string) {
	c.defaults.ID = id
}

func (c *CodeBuilder) SetTTL(ttl time.Duration) {
	c.defaults.TTL = ttl
}

func (c *CodeBuilder) Build() *Code {
	result := c.defaults
	if c.overwrite != nil {
		c.overwrite(&result)
	}
	return &result
}
