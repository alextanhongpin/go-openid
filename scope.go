package oidc

type Scope struct {
	OpenID        bool
	Profile       bool
	Email         bool
	Address       bool
	Phone         bool
	OfflineAccess bool
}
