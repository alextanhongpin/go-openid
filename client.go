package oidc

import jwt "github.com/dgrijalva/jwt-go"

//go:generate gencodec -type Client -out gen_client.go

// ClientToken represents the access token that is provided to the client
// during registration.
type ClientToken struct {
	*jwt.StandardClaims
	ClientID string
}

// Client represents the OIDC Client Metadata.
type Client struct {
	ApplicationType              string       `json:"application_type,omitempty"`
	ClientName                   string       `json:"client_name,omitempty"`
	ClientURI                    string       `json:"client_uri,omitempty"`
	Contacts                     []string     `json:"contacts,omitempty"`
	DefaultAcrValues             string       `json:"default_acr_values,omitempty"`
	DefaultMaxAge                int64        `json:"default_maxa_age,omitempty"`
	GrantTypes                   []string     `json:"grant_types,omitempty"`
	IDTokenEncryptedResponseAlg  string       `json:"id_token_encrypted_response_alg,omitempty"`
	IDTokenEncryptedResponseEnc  string       `json:"id_token_encryption_response_enc,omitempty"`
	IDTokenSignedResponseAlg     string       `json:"id_token_signed_response_alg,omitempty"`
	InitiateLoginURI             string       `json:"initiate_login_uri,omitempty"`
	Jwks                         string       `json:"jwks,omitempty"`
	JwksURI                      string       `json:"jwks_uri,omitempty"`
	LogoURI                      string       `json:"logo_uri,omitempty"`
	PolicyURI                    string       `json:"policy_uri,omitempty"`
	RedirectURIs                 RedirectURIs `json:"redirect_uris,omitempty"`
	RequestObjectEncryptionAlg   string       `json:"request_object_encryption_alg,omitempty"`
	RequestObjectEncryptionEnc   string       `json:"request_object_encryption_enc,omitempty"`
	RequestObjectSigningAlg      string       `json:"request_object_signing_alg,omitempty"`
	RequestURIs                  []string     `json:"request_uris,omitempty"`
	RequireAuthTime              int64        `json:"require_auth_time,omitempty"`
	ResponseTypes                []string     `json:"response_types,omitempty"`
	SectorIdentifierURI          string       `json:"sector_identifier_uri,omitempty"`
	SubjectType                  string       `json:"subject_type,omitempty"`
	TokenEndpointAuthMethod      string       `json:"token_endpoint_auth_method,omitempty"`
	TokenEndpointAuthSigningAlg  string       `json:"token_endpoint_auth_signing_alg,omitempty"`
	TosURI                       string       `json:"tos_uri,omitempty"`
	UserinfoEncryptedResponseAlg string       `json:"userinfo_encrypted_response_alg,omitempty"`
	UserinfoEncryptedResponseEnc string       `json:"userinfo_encrypted_response_enc,omitempty"`
	UserinfoSignedResponseAlg    string       `json:"userinfo_signed_response_alg,omitempty"`
	ClientID                     string       `json:"client_id,omitempty"`
	ClientIDIssuedAt             int64        `json:"client_id_issued_at,omitempty"`
	ClientSecret                 string       `json:"client_secret,omitempty"`
	ClientSecretExpiresAt        int64        `json:"client_secret_expires_at,omitempty"`
	RegistrationAccessToken      string       `json:"registration_access_token,omitempty"`
	RegistrationClientURI        string       `json:"registration_client_uri,omitempty"`
}

// NewClient returns a new client with the given name.
func NewClient(name string) *Client {
	return &Client{ClientName: name}
}

func (c *Client) Copy() *Client {
	copy := new(Client)
	*copy = *c
	return copy
}

// -- redirect uris

// RedirectURIs represents a slice of valid redirect uris.
type RedirectURIs []string

// Contains checks if the redirect uri is present in the slice.
func (r RedirectURIs) Contains(uri string) bool {
	for _, u := range r {
		if u == uri {
			return true
		}
	}
	return false
}
