package database

// ClientRepository represents the interface for the client repository.
type ClientRepository interface {
	Get(name string) (*oidc.Client, bool)
	GetByID(id string) *oidc.Client
	GetByIDAndSecret(id, secret string) *oidc.Client
	Put(id string, client *oidc.Client)
	Delete(name string)
}
