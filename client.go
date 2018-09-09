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
	InvalidClientMetadata ClientErrorCode = iota
	InvalidRedirectURI
)

var clientErrorDescriptions = map[ClientErrorCode]string{
	InvalidClientMetadata: "the value of one of the client metadata fields is invalid and the server has rejected this request",
	InvalidRedirectURI:    "the value of one or more redirect uris is invalid",
}

var clientErrorCodes = map[ClientErrorCode]string{
	InvalidClientMetadata: "invalid_client_metadata",
	InvalidRedirectURI:    "invalid_redirect_uri",
}

// String fulfills the Stringer method.
func (c ClientErrorCode) String() string {
	return clientErrorCodes[c]
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
		State:       "",
		URI:         "",
	}
}

// Client represents both private and public metadata of the client.
type Client struct {
	*ClientPrivate
	*ClientPublic
}

// NewClient returns a client metadata given the public client metadata
func NewClient(req *ClientRegistrationRequest) *Client {
	return &Client{
		ClientPublic: req,
		ClientPrivate: &ClientRegistrationResponse{
			ClientID:                "fake client id",
			ClientIDIssuedAt:        0,
			ClientSecret:            "fake client secret",
			ClientSecretExpiresAt:   0,
			RegistrationAccessToken: "",
			RegistrationClientURI:   "",
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
	ClientID                string `json:"client_id,omitempty"`
	ClientIDIssuedAt        int64  `json:"client_id_issued_at,omitempty"`
	ClientSecret            string `json:"client_secret,omitempty"`
	ClientSecretExpiresAt   int64  `json:"client_secret_expires_at,omitempty"`
	RegistrationAccessToken string `json:"registration_access_token,omitempty"`
	RegistrationClientURI   string `json:"registration_client_uri,omitempty"`
}
