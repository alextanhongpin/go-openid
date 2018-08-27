package oidc

const ClientRegistrationEndpoint = "/connect/register"

type Client struct {
	*ClientRegistrationRequest
	*ClientRegistrationResponse
}

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

type ClientRegistrationRequest struct {
	RedirectURIs                 []string `json:"redirect_uris"`
	ResponseTypes                []string `json:"response_types"`
	GrantTypes                   []string `json:"grant_types"`
	ApplicationType              string   `json:"application_type"`
	Contacts                     []string `json:"contacts"`
	ClientName                   string   `json:"client_name"`
	LogoURI                      string
	ClientURI                    string
	PolicyURI                    string
	TosURI                       string
	JwksURI                      string
	Jwks                         string
	SectorIdentifierURI          string
	SubjectType                  string
	IDTokenSignedResponseAlg     string
	IDTokenEncryptedResponseAlg  string
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

func (c *ClientRegistrationRequest) Validate() error {
	return nil
}

type ClientRegistrationResponse struct {
	ClientID                string
	ClientSecret            string
	RegistrationAccessToken string
	RegistrationClientURI   string
	ClientIDIssuedAt        int64
	ClientSecretExpiresAt   int64
}

type ClientErrorResponse struct {
	Error            string
	ErrorDescription string
}
