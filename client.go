package oidc

const ClientRegistrationEndpoint = "/connect/register"

// Client represents both private and public metadata of the client.
type Client struct {
	*ClientRegistrationRequest
	*ClientRegistrationResponse
}

// NewClient returns a client metadata given the public client metadata
func NewClient(req *ClientRegistrationRequest) *Client {
	return &Client{
		ClientRegistrationRequest: req,
		ClientRegistrationResponse: &ClientRegistrationResponse{
			ClientID:                "fake client id",
			ClientSecret:            "fake client secret",
			RegistrationAccessToken: "",
			RegistrationClientURI:   "",
			ClientIDIssuedAt:        0,
			ClientSecretExpiresAt:   0,
		},
	}
}

// ClientRegistrationRequest represents the client registration request.
type ClientRegistrationRequest struct {
	RedirectURIs                 []string `json:"redirect_uris"`
	ResponseTypes                []string `json:"response_types"`
	GrantTypes                   []string `json:"grant_types"`
	ApplicationType              string   `json:"application_type"`
	Contacts                     []string `json:"contacts"`
	ClientName                   string   `json:"client_name"`
	LogoURI                      string   `json:"logo_uri"`
	ClientURI                    string   `json:"client_uri"`
	PolicyURI                    string   `json:"policy_uri"`
	TosURI                       string   `json:"tos_uri"`
	JwksURI                      string   `json:"jwks_uri"`
	Jwks                         string   `json:"jwks"`
	SectorIdentifierURI          string   `json:"sector_identifier_uri"`
	SubjectType                  string   `json:"subject_type"`
	IDTokenSignedResponseAlg     string   `json:"id_token_signed_response_alg"`
	IDTokenEncryptedResponseAlg  string   `json:"id_token_encrypted_response_alg"`
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
	// Check the redirect uri
	// return ErrInvalidRedirectURI

	// Check the client metadata
	// return ErrInvalidClientMetadata
	return nil
}

// ClientRegistrationResponse represents the response payload of the client.
type ClientRegistrationResponse struct {
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
