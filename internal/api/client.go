package main

import (
	"context"
	"time"
)

// Client represents the openid Client Metadata.
type Client struct {
	ApplicationType              string   `json:"application_type,omitempty"`
	ClientName                   string   `json:"client_name,omitempty"`
	ClientURI                    string   `json:"client_uri,omitempty"`
	Contacts                     []string `json:"contacts,omitempty"`
	DefaultAcrValues             string   `json:"default_acr_values,omitempty"`
	DefaultMaxAge                int64    `json:"default_maxa_age,omitempty"`
	GrantTypes                   []string `json:"grant_types,omitempty"`
	IDTokenEncryptedResponseAlg  string   `json:"id_token_encrypted_response_alg,omitempty"`
	IDTokenEncryptedResponseEnc  string   `json:"id_token_encryption_response_enc,omitempty"`
	IDTokenSignedResponseAlg     string   `json:"id_token_signed_response_alg,omitempty"`
	InitiateLoginURI             string   `json:"initiate_login_uri,omitempty"`
	Jwks                         string   `json:"jwks,omitempty"`
	JwksURI                      string   `json:"jwks_uri,omitempty"`
	LogoURI                      string   `json:"logo_uri,omitempty"`
	PolicyURI                    string   `json:"policy_uri,omitempty"`
	RedirectURIs                 []string `json:"redirect_uris,omitempty"`
	RequestObjectEncryptionAlg   string   `json:"request_object_encryption_alg,omitempty"`
	RequestObjectEncryptionEnc   string   `json:"request_object_encryption_enc,omitempty"`
	RequestObjectSigningAlg      string   `json:"request_object_signing_alg,omitempty"`
	RequestURIs                  []string `json:"request_uris,omitempty"`
	RequireAuthTime              int64    `json:"require_auth_time,omitempty"`
	ResponseTypes                []string `json:"response_types,omitempty"`
	SectorIdentifierURI          string   `json:"sector_identifier_uri,omitempty"`
	SubjectType                  string   `json:"subject_type,omitempty"`
	TokenEndpointAuthMethod      string   `json:"token_endpoint_auth_method,omitempty"`
	TokenEndpointAuthSigningAlg  string   `json:"token_endpoint_auth_signing_alg,omitempty"`
	TosURI                       string   `json:"tos_uri,omitempty"`
	UserinfoEncryptedResponseAlg string   `json:"userinfo_encrypted_response_alg,omitempty"`
	UserinfoEncryptedResponseEnc string   `json:"userinfo_encrypted_response_enc,omitempty"`
	UserinfoSignedResponseAlg    string   `json:"userinfo_signed_response_alg,omitempty"`
	ClientID                     string   `json:"client_id,omitempty"`
	ClientIDIssuedAt             int64    `json:"client_id_issued_at,omitempty"`
	ClientSecret                 string   `json:"client_secret,omitempty"`
	ClientSecretExpiresAt        int64    `json:"client_secret_expires_at,omitempty"`
	RegistrationAccessToken      string   `json:"registration_access_token,omitempty"`
	RegistrationClientURI        string   `json:"registration_client_uri,omitempty"`
}

func (c *Client) HasRedirectURI(uri string) bool {
	return URIs(c.RedirectURIs).Contains(uri)
}

// NewClient returns a new client with default values.
func NewClient() *Client {
	return &Client{
		ApplicationType:              "web",
		GrantTypes:                   []string{"authorization_code"},
		RequestObjectEncryptionEnc:   "A128CBC-HS256",
		ResponseTypes:                []string{"code"},
		UserinfoEncryptedResponseEnc: "A128CBC-HS256",
	}
}

type StringGenerator func(n int) (string, error)

func RegisterClient(ctx context.Context, stringGenerator StringGenerator, req *Client) (*Client, error) {
	// TODO: Validate URI.
	now, ok := ctx.Value(ContextKeyTimestamp).(time.Time)
	if !ok {
		now = time.Now().UTC()
	}

	var err error
	req.ClientSecret, err = stringGenerator(32)
	if err != nil {
		return nil, err
	}
	req.ClientID, err = stringGenerator(32)
	if err != nil {
		return nil, err
	}
	// req.RegistrationAccessToken = ""
	// req.RegistrationClientURI = ""
	req.ClientIDIssuedAt = now.Unix()
	req.ClientSecretExpiresAt = 0
	return req, nil
}

func ReadClient(repo ClientRepository, clientID string) (*Client, error) {
	return repo.GetClientByClientID(clientID)
}
