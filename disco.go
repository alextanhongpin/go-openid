package oidc

// Discovery represents the OpenID Discovery protocol.
type Discovery struct {
	Resource string `json:"resource,omitempty"`
	Host     string `json:"host,omitempty"`
	Rel      string `json:"rel,omitempty"`
}
