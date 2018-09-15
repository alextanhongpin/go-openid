package database

import "github.com/alextanhongpin/go-openid"

// ClientRepo represents the interface for the client repository.
type ClientRepo interface {
	Get(name string) (*oidc.Client, bool)
	GetByID(id string) *oidc.Client
	GetByIDAndSecret(id, secret string) *oidc.Client
	Put(id string, client *oidc.Client)
	Delete(name string)
}

// CodeRepo represents the operations for the code repository.
type CodeRepo interface {
	Get(id string) (*oidc.Code, bool)
	Put(id string, code *oidc.Code)
	Delete(id string)
}

// UserRepo represents the operations for the user repository.
type UserRepo interface {
	Get(id string) (*oidc.IDToken, bool)
	Put(id string, user *oidc.IDToken)
	Delete(id string)
}
