package client

import "errors"

const ClientRegistrationEndpoint = "/connect/register"

var (
	ErrInvalidRedirectURI    = errors.New("invalid redirect uri")
	ErrInvalidClientMetadata = errors.New("invalid client metadata")
)

type Client struct {
	RedirectURIs                 string
	ResponseTypes                []string
	GrantTypes                   []string
	ApplicationType              string
	Contacts                     []string
	ClientName                   string
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

func New() *Client {
	return &Client{}
}

type ClientRegistrationRequest struct{}
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

type ClientService interface {
	Register(ClientRegistrationRequest) (*ClientRegistrationResponse, error)
}

type clientService struct {
}

func (c *clientService) Register(req ClientRegistrationRequest) (*ClientRegistrationResponse, error) {
	return nil, nil
}
