package oidc

// Discovery represents the OpenID Discovery protocol.
type Discovery struct {
	Host     string `json:"host,omitempty"`
	Rel      string `json:"rel,omitempty"`
	Resource string `json:"resource,omitempty"`
}
