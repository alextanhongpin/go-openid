package main

type (
	ClientRepository interface {
		GetClientByClientID(clientID string) (*Client, error)
		GetClientByCredentials(clientID, clientSecret string) (*Client, error)
	}
	CodeRepository interface {
		GetCodeByID(id string) (*Code, error)
		Create(code *Code) error
	}
)
