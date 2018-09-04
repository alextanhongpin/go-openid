package oidc

import (
	"strings"

	"github.com/asaskevich/govalidator"
)

// ClientRegistrationEndpoint represents the client registration endpoint.
const ClientRegistrationEndpoint = "/connect/register"

// ClientErrorCode represents the error code returned by client.
type ClientErrorCode int

const (
	InvalidRedirectURI ClientErrorCode = iota
	InvalidClientMetadata
)

var clientErrorDescriptions = map[ClientErrorCode]string{
	InvalidRedirectURI:    "the value of one or more redirect uris is invalid",
	InvalidClientMetadata: "the value of one of the client metadata fields is invalid and the server has rejected this request",
}

// String fulfills the Stringer method.
func (c ClientErrorCode) String() string {
	return [...]string{"invalid_redirect_uri", "invalid_client_metadata"}[c]
}

// Description returns the general client error description.
func (c ClientErrorCode) Description() string {
	return clientErrorDescriptions[c]
}

// JSON returns the error json.
func (c ClientErrorCode) JSON() *ErrorJSON {
	return &ErrorJSON{
		Code:        c.String(),
		Description: c.Description(),
		URI:         "",
		State:       "",
	}
}

// Client represents both private and public metadata of the client.
type Client struct {
	*ClientPublic
	*ClientPrivate
}

// NewClient returns a client metadata given the public client metadata
func NewClient(req *ClientRegistrationRequest) *Client {
	return &Client{
		ClientPublic: req,
		ClientPrivate: &ClientRegistrationResponse{
			ClientID:                "fake client id",
			ClientSecret:            "fake client secret",
			RegistrationAccessToken: "",
			RegistrationClientURI:   "",
			ClientIDIssuedAt:        0,
			ClientSecretExpiresAt:   0,
		},
	}
}

// RedirectURIs represents a slice of valid redirect uris.
type RedirectURIs []string

// Contains checks if the redirect uri is present in the slice.
func (r RedirectURIs) Contains(uri string) bool {
	for _, u := range r {
		if strings.Compare(u, uri) == 0 {
			return true
		}
	}
	return false
}

// ClientRegistrationRequest represents the client registration request.
type ClientRegistrationRequest = ClientPublic

// ClientPublic represents fields that are public
type ClientPublic struct {
	RedirectURIs                 RedirectURIs `json:"redirect_uris,omitempty"`
	ResponseTypes                []string     `json:"response_types,omitempty"`
	GrantTypes                   []string     `json:"grant_types,omitempty"`
	ApplicationType              string       `json:"application_type,omitempty"`
	Contacts                     []string     `json:"contacts,omitempty"`
	ClientName                   string       `json:"client_name,omitempty"`
	LogoURI                      string       `json:"logo_uri,omitempty"`
	ClientURI                    string       `json:"client_uri,omitempty"`
	PolicyURI                    string       `json:"policy_uri,omitempty"`
	TosURI                       string       `json:"tos_uri,omitempty"`
	JwksURI                      string       `json:"jwks_uri,omitempty"`
	Jwks                         string       `json:"jwks,omitempty"`
	SectorIdentifierURI          string       `json:"sector_identifier_uri,omitempty"`
	SubjectType                  string       `json:"subject_type,omitempty"`
	IDTokenSignedResponseAlg     string       `json:"id_token_signed_response_alg,omitempty"`
	IDTokenEncryptedResponseAlg  string       `json:"id_token_encrypted_response_alg,omitempty"`
	IDTokenEncryptedResponseEnc  string
	UserinfoSignedResponseAlg    string
	UserinfoEncryptedResponseAlg string
	UserinfoEncryptedResponseEnc string
	RequestObjectSigningAlg      string
	RequestObjectEncryptionAlg   string
	RequestObjectEncryptionEnc   string
	TokenEndpointAuthMethod      string
	TokenEndpointAuthSigningAlg  string
	DefaultMaxAge                int64
	RequireAuthTime              int64
	DefaultAcrValues             string
	InitiateLoginURI             string
	RequestURIs                  []string
}

// Validate performs a simple validation on the client payload request.
func (c *ClientRegistrationRequest) Validate() error {
	for _, u := range c.RedirectURIs {
		if !govalidator.IsURL(u) {
			return InvalidRedirectURI.JSON()
		}
	}
	// Check the redirect uri
	// return ErrInvalidRedirectURI

	// Check the client metadata
	// return ErrInvalidClientMetadata
	return nil
}

// ClientRegistrationResponse represents the response payload of the client.
type ClientRegistrationResponse = ClientPrivate

// ClientPrivate represents fields that are private.
type ClientPrivate struct {
	ClientID                string `json:"client_id"`
	ClientSecret            string `json:"client_secret"`
	RegistrationAccessToken string `json:"registration_access_token"`
	RegistrationClientURI   string `json:"registration_client_uri"`
	ClientIDIssuedAt        int64  `json:"client_id_issued_at"`
	ClientSecretExpiresAt   int64  `json:"client_secret_expires_at"`
}

// ClientErrorResponse represents the error response of the client.
type ClientErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
