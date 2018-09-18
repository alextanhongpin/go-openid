package repository

import "github.com/alextanhongpin/go-openid"

// Code represents the operations for the code repository.
type Code interface {
	Get(id string) (*oidc.Code, bool)
	Put(id string, code *oidc.Code)
	Delete(id string)
}
