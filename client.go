package oidc

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
)

//go:generate gencodec -type Client -out gen_client.go

// ClientRegistrationEndpoint represents the client registration endpoint.
const ClientRegistrationEndpoint = "/connect/register"

// Client represents fields that are public.
type Client struct {
	// Public
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
	// Private read-only fields.
	ClientID                string `json:"-"`
	ClientIDIssuedAt        int64  `json:"-"`
	ClientSecret            string `json:"-"`
	ClientSecretExpiresAt   int64  `json:"-"`
	RegistrationAccessToken string `json:"-"`
}

type ClientRegistrationRequest = Client

// Validate performs a simple validation on the client payload request.
func (c *ClientRegistrationRequest) Validate() error {
	if u, err := validateURIs(c.RedirectURIs...); err != nil {
		desc := fmt.Sprintf(`%s is not a valid uri format`, u)
		return ErrInvalidRedirectURI.WithDescription(desc)
	}

	// Validate Response type
	// Validate grant types
	// Validate application type
	// Validate contacts
	// Validate client name
	// Validate subject type: pairwise and public
	// Validase ...

	// Check the client metadata
	// return ErrInvalidClientMetadata

	uris := []string{
		c.ClientURI,
		c.LogoURI,
		c.PolicyURI,
		c.TosURI,
		c.JwksURI,
		c.SectorIdentifierURI,
		c.InitiateLoginURI,
	}
	uris = append(uris, c.RedirectURIs...)
	if _, err := validateURIs(uris...); err != nil {
		return err
	}
	return nil
}

// ClientRegistrationResponse represents the response payload of the client.
type ClientRegistrationResponse struct {
	ClientID                string `json:"client_id,omitempty"`
	ClientIDIssuedAt        int64  `json:"client_id_issued_at,omitempty"`
	ClientSecret            string `json:"client_secret,omitempty"`
	ClientSecretExpiresAt   int64  `json:"client_secret_expires_at,omitempty"`
	RegistrationAccessToken string `json:"registration_access_token,omitempty"`
	RegistrationClientURI   string `json:"registration_client_uri,omitempty"`
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

type GrantTypes []string

var grantypesmap = map[string]struct{}{
	"authorization_code": struct{}{},
	"implicit":           struct{}{},
	"refresh_token":      struct{}{},
}

func (g GrantTypes) Validate() bool {
	for _, v := range g {
		if _, ok := grantypesmap[v]; !ok {
			return false
		}
	}
	return true
}

func validateApplicationType(in string) error {
	// OPTIONAL. Kind of the application. The default, if omitted, is
	// web. The defined values are native or web. Web Clients using the
	// OAuth Implicit Grant Type MUST only register URLs using the https
	// scheme as redirect_uris; they MUST NOT use localhost as the
	// hostname. Native Clients MUST only register redirect_uris using
	// custom URI schemes or URLs using the http: scheme with localhost
	// as the hostname. Authorization Servers MAY place additional
	// constraints on Native Clients. Authorization Servers MAY reject
	// Redirection URI values using the http scheme, other than the
	// localhost case for Native Clients. The Authorization Server MUST
	// verify that all the registered redirect_uris conform to these
	// constraints. This prevents sharing a Client ID across different
	// types of Clients.
	// if in != "web" || in != "native" {
	// }
	return nil
}

func validateContacts(in []string) error {
	for _, v := range in {
		if !govalidator.IsEmail(v) {
			return fmt.Errorf("invalid email format for %s", v)
		}
	}
	return nil
}

func validateClientResponseType(in string) error {
	return nil
}

func validateURIs(uris ...string) (string, error) {
	for _, u := range uris {
		if u == "" {
			// Allow optional strings
			continue
		}
		if !govalidator.IsURL(u) {
			return u, errors.New("invalid url format")
		}
	}
	return "", nil
}
