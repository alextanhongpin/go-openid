package model

// Core represents the core model for the OpenID Connect Specification.
type Core interface {
	ValidateClient(clientID, redirectURI string) error
	NewCode() string
}
