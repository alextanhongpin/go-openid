package main

type ClientRepository interface {
	GetClientByClientID(clientID string) (*Client, error)
	GetClientByCredentials(clientID, clientSecret string) (*Client, error)
}
type CodeRepository interface {
	GetCodeByID(id string) (*Code, error)
}
