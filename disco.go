package oidc

// Discovery represents the OpenID Discovery protocol.
type Discovery struct {
	Resource string
	Host     string
	Rel      string
}
